package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/autoscale"
	"github.com/jrasell/sherpa/pkg/client"
	policyBackend "github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/policy/backend/consul"
	policyMemory "github.com/jrasell/sherpa/pkg/policy/backend/memory"
	"github.com/jrasell/sherpa/pkg/policy/backend/nomadmeta"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/jrasell/sherpa/pkg/server/cluster"
	"github.com/jrasell/sherpa/pkg/server/router"
	clusterBackend "github.com/jrasell/sherpa/pkg/state/cluster"
	clusterConsul "github.com/jrasell/sherpa/pkg/state/cluster/consul"
	clusterMemory "github.com/jrasell/sherpa/pkg/state/cluster/memory"
	stateBackend "github.com/jrasell/sherpa/pkg/state/scale"
	stateConsul "github.com/jrasell/sherpa/pkg/state/scale/consul"
	stateMemory "github.com/jrasell/sherpa/pkg/state/scale/memory"
	"github.com/jrasell/sherpa/pkg/watcher"
	"github.com/jrasell/sherpa/pkg/watcher/deployment"
	"github.com/jrasell/sherpa/pkg/watcher/job"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	addr   string
	cfg    *Config
	logger zerolog.Logger

	policyBackend  policyBackend.PolicyBackend
	stateBackend   stateBackend.Backend
	scaleBackend   scale.Scale
	clusterBackend clusterBackend.Backend

	// deploymentWatcher is used to watch deployments in order to update internal tracking.
	deploymentWatcher watcher.Watcher

	// nomadMetaWatcher is used to watch Nomad jobs in order to update policies based off the Nomad
	// meta stanzas.
	nomadMetaWatcher watcher.Watcher

	// nomadMetaProcessor is the processor which is used to handle job updates and decide if their
	// scaling meta policy has changed and should be reflected in storage.
	nomadMetaProcessor *nomadmeta.Processor

	clusterMember *cluster.Member

	nomad     *api.Client
	autoScale *autoscale.AutoScale
	telemetry *metrics.InmemSink

	http.Server
	routes *routes

	// gcIsRunning is used to track whether this Sherpa server is currently running the garbage
	// collection loop.
	gcIsRunning bool

	// stopChan is used to synchronise stopping the HTTP server services and any handlers which it
	// maintains operationally.
	stopChan chan struct{}
}

func New(l zerolog.Logger, cfg *Config) *HTTPServer {
	return &HTTPServer{
		addr:     fmt.Sprintf("%s:%d", cfg.Server.Bind, cfg.Server.Port),
		cfg:      cfg,
		logger:   l,
		routes:   &routes{},
		stopChan: make(chan struct{}),
	}
}

func (h *HTTPServer) Start() error {
	h.logger.Info().Str("addr", h.addr).Msg("starting HTTP server")
	h.logServerConfig()

	if err := h.setup(); err != nil {
		return err
	}

	go h.leaderUpdateHandler()

	// Start the deployment watcher, using the scale deployment channel for updates.
	go h.deploymentWatcher.Run(h.scaleBackend.GetDeploymentChannel())

	// If the operator has configured the Nomad meta policy engine, we should start the processes
	// which watch and handle updates.
	if h.cfg.Server.NomadMetaPolicyEngine && h.nomadMetaWatcher != nil {
		go h.nomadMetaProcessor.Run()
		go h.nomadMetaWatcher.Run(h.nomadMetaProcessor.GetUpdateChannel())
	}

	h.handleSignals()
	return nil
}

func (h *HTTPServer) logServerConfig() {
	h.logger.Info().
		Object("server", h.cfg.Server).
		Object("tls", h.cfg.TLS).
		Object("telemetry", h.cfg.Telemetry).
		Object("cluster", h.cfg.Cluster).
		Msg("Sherpa server configuration")
}

func (h *HTTPServer) setup() error {
	if err := h.setupNomadClient(); err != nil {
		return err
	}

	// Setup telemetry based on the config passed by the operator.
	if err := h.setupTelemetry(); err != nil {
		return errors.Wrap(err, "failed to setup telemetry handler")
	}

	h.setupStoredBackends()

	h.setupScaler()
	go h.scaleBackend.RunDeploymentUpdateHandler()

	h.setupDeploymentWatcher()

	mem, err := cluster.NewMember(h.logger, h.clusterBackend, h.addr, h.cfg.Cluster.Addr, h.cfg.Cluster.Name)
	if err != nil {
		return err
	}
	h.clusterMember = mem

	go h.clusterMember.RunLeadershipLoop()

	// If the server has been set to enable the internal autoscaler, set this up. We should not
	// start the handler here; it is the responsibility of the leadership update handler to perform
	// this task.
	if h.cfg.Server.InternalAutoScaler {
		if err := h.setupAutoScaling(); err != nil {
			return errors.Wrap(err, "failed to setup internal autoscaler")
		}
	}

	initialRoutes := h.setupRoutes()

	r := router.WithRoutes(h.logger, *initialRoutes)
	http.Handle("/", middlewareLogger(r, h.logger))

	// Run the TLS setup process so that if the user has configured a TLS certificate pair the
	// server uses these.
	if err := h.setupTLS(); err != nil {
		return err
	}

	// Once we have the TLS config in place, we can setup the listener which uses the TLS setup to
	// correctly start the listener.
	ln := h.setupListener()
	if ln == nil {
		return errors.New("failed to setup HTTP server, listener is nil")
	}
	h.logger.Info().Str("addr", h.addr).Msg("HTTP server successfully listening")

	go func() {
		err := h.Serve(ln)
		h.logger.Info().Msgf("HTTP server has been shutdown: %v", err)
	}()

	return nil
}

func (h *HTTPServer) setupTLS() error {
	if h.cfg.TLS.CertPath != "" && h.cfg.TLS.CertKeyPath != "" {
		h.logger.Debug().Msg("setting up server TLS")

		cert, err := tls.LoadX509KeyPair(h.cfg.TLS.CertPath, h.cfg.TLS.CertKeyPath)
		if err != nil {
			return errors.Wrap(err, "failed to load certificate cert/key pair")
		}
		h.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}
	return nil
}

func (h *HTTPServer) setupStoredBackends() {

	// Setup the standard backends based on the operators storage type.
	if h.cfg.Server.ConsulStorageBackend {
		h.logger.Debug().Msg("setting up Consul storage backend")
		h.stateBackend = stateConsul.NewStateBackend(h.logger, h.cfg.Server.ConsulStorageBackendPath)
		h.clusterBackend = clusterConsul.NewStateBackend(h.logger, h.cfg.Server.ConsulStorageBackendPath)
	} else {
		h.logger.Debug().Msg("setting up in-memory storage backend")
		h.stateBackend = stateMemory.NewStateBackend()
		h.clusterBackend = clusterMemory.NewStateBackend()
	}
	h.setupPolicyBackend()
}

func (h *HTTPServer) setupPolicyBackend() {
	h.logger.Debug().Msg("setting up policy backend")

	if h.cfg.Server.NomadMetaPolicyEngine {
		h.nomadMetaWatcher = job.NewWatcher(h.logger, h.nomad)
		h.policyBackend, h.nomadMetaProcessor = nomadmeta.NewJobScalingPolicies(h.logger, h.nomad)
		return
	}

	if h.cfg.Server.ConsulStorageBackend {
		h.policyBackend = consul.NewConsulPolicyBackend(h.logger, h.cfg.Server.ConsulStorageBackendPath)
		return
	}
	h.policyBackend = policyMemory.NewJobScalingPolicies()
}

func (h *HTTPServer) setupNomadClient() error {
	h.logger.Debug().Msg("setting up Nomad client")

	nc, err := client.NewNomadClient()
	if err != nil {
		return err
	}
	h.nomad = nc

	return nil
}

func (h *HTTPServer) setupAutoScaling() error {
	h.logger.Debug().Msg("setting up Sherpa internal auto-scaling engine")
	autoscaleCfg := &autoscale.Config{
		StrictChecking:  h.cfg.Server.StrictPolicyChecking,
		ScalingInterval: h.cfg.Server.InternalAutoScalerEvalPeriod,
		ScalingThreads:  h.cfg.Server.InternalAutoScalerNumThreads,
	}

	as, err := autoscale.NewAutoScaleServer(h.logger, h.nomad, h.policyBackend, h.scaleBackend, autoscaleCfg)
	if err != nil {
		return err
	}
	h.autoScale = as

	return nil
}

func (h *HTTPServer) setupScaler() {
	h.scaleBackend = scale.NewScaler(h.nomad, h.logger, h.stateBackend, h.cfg.Server.StrictPolicyChecking)
}

func (h *HTTPServer) setupDeploymentWatcher() {
	h.deploymentWatcher = deployment.New(h.logger, h.nomad)
}

func (h *HTTPServer) setupListener() net.Listener {
	var (
		err error
		ln  net.Listener
	)

	if h.TLSConfig != nil {
		ln, err = tls.Listen("tcp", h.addr, h.TLSConfig)
	} else {
		ln, err = net.Listen("tcp", h.addr)
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to setup server HTTP listener")
	}
	return ln
}

// leaderUpdateHandler is the server process which monitors for messages from the leadership
// process. Changes in the leadership status of a Sherpa server means the server will need to start
// or stop a number of sub-process.
func (h *HTTPServer) leaderUpdateHandler() {
	for {
		select {
		case <-h.stopChan:
			h.logger.Info().Msg("shutting down server leader update handler")
			return
		case msg := <-h.clusterMember.UpdateChan:
			h.logger.Debug().Str("leadership-msg", msg.Msg).Msg("server received leader update message")
			h.handleLeaderUpdateMsg(msg.IsLeader)
		}
	}
}

// handleLeaderUpdateMsg is responsible for acting on a leadership message and performing the
// start/stop actions as a result.
func (h *HTTPServer) handleLeaderUpdateMsg(isLeader bool) {
	switch isLeader {
	case true:
		if h.autoScale != nil && !h.autoScale.IsRunning() {
			go h.autoScale.Run()
		}
		if !h.gcIsRunning {
			go h.runGarbageCollectionLoop()
		}
	default:
		if h.autoScale != nil && h.autoScale.IsRunning() {
			h.autoScale.Stop()
		}
		if h.gcIsRunning {
			h.stopChan <- struct{}{}
		}
	}
}

// Stop is used to synchronise the shutdown of background tasks before the server exits.
func (h *HTTPServer) Stop() error {
	h.logger.Info().Msg("gracefully shutting down HTTP server and sub-processes")

	// If the autoscaler is running, stop this. It is important that a Sherpa server is given time
	// to exit cleanly as this call can take a number of seconds to complete while we gracefully
	// wait for all in-flight worker threads to finish.
	if h.autoScale != nil && h.autoScale.IsRunning() {
		h.autoScale.Stop()
	}

	// Stop the leadership loop and remove any stored leadership information. It is not important
	// that this happens cleanly, but preferred.
	h.clusterMember.ClearLeadership()

	// Send a signal to the HTTPServer stopChan instructing sub-process to stop.
	close(h.stopChan)

	// When calling shutdown, the process will wait for all active connections to finish. This
	// protects against interrupting scaling events triggered via the API.
	return h.Shutdown(context.Background())
}

// handleSignals is responsible for blocking on receiving OS signals, then handling signals as
// required by the Sherpa server.
func (h *HTTPServer) handleSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-signalCh:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				if err := h.Stop(); err != nil {
					h.logger.Error().Err(err).Msg("failed to cleanly shutdown server and sub-processes")
				}
				h.logger.Info().Msg("successfully shutdown server and sub-processes")
				return
			default:
				panic(fmt.Sprintf("unsupported signal: %v", sig))
			}
		}
	}
}
