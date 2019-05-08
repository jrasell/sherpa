package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/autoscale"
	"github.com/jrasell/sherpa/pkg/client"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/policy/backend/consul"
	"github.com/jrasell/sherpa/pkg/policy/backend/memory"
	"github.com/jrasell/sherpa/pkg/server/router"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	addr   string
	cfg    *Config
	logger zerolog.Logger

	nomad         *api.Client
	policyBackend backend.PolicyBackend
	autoScale     *autoscale.AutoScale

	http.Server
	routes *routes
}

func New(l zerolog.Logger, cfg *Config) *HTTPServer {
	return &HTTPServer{
		addr:   fmt.Sprintf("%s:%d", cfg.Server.Bind, cfg.Server.Port),
		cfg:    cfg,
		logger: l,
		routes: &routes{},
	}
}

func (h *HTTPServer) Start() error {
	h.logger.Info().Str("addr", h.addr).Msg("starting HTTP server")
	h.logger.Info().Object("config", h.cfg.Server).Msg("Sherpa server configuration")
	if err := h.setup(); err != nil {
		return err
	}

	h.handleSignals(context.Background())
	return nil
}

func (h *HTTPServer) setup() error {
	if err := h.setupNomadClient(); err != nil {
		return err
	}

	h.setupPolicyBackend()

	// If the server has been set to enable the internal autoscaler, set this up and start the
	// process running.
	if h.cfg.Server.InternalAutoScaler {
		h.setupAutoScaling()
	}

	initialRoutes := h.setupRoutes()

	r := router.WithRoutes(h.logger, *initialRoutes)
	http.Handle("/", middlewareLogger(r, h.logger))

	ln := h.listenWithRetry()

	go func() {
		err := h.Serve(ln)
		h.logger.Info().Err(err).Msg("HTTP server has been shutdown")
	}()

	return nil
}

func (h *HTTPServer) setupPolicyBackend() {
	if h.cfg.Server.ConsulStorageBackend {
		h.logger.Debug().Msg("setting up Consul storage backend")
		h.policyBackend = consul.NewConsulPolicyBackend(h.logger, h.cfg.Server.ConsulStorageBackendPath)
	} else {
		h.logger.Debug().Msg("setting up in-memory storage backend")
		h.policyBackend = memory.NewJobScalingPolicies()
	}
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

func (h *HTTPServer) setupAutoScaling() {
	h.logger.Debug().Msg("setting up Sherpa internal auto-scaling engine")
	autoscaleCfg := &autoscale.Config{StrictChecking: h.cfg.Server.StrictPolicyChecking, Scaling: h.cfg.AutoScale}
	h.autoScale = autoscale.NewAutoScaleServer(h.logger, h.nomad, h.policyBackend, autoscaleCfg)
	go h.autoScale.Run()
}

func (h *HTTPServer) listenWithRetry() net.Listener {
	var (
		err error
		ln  net.Listener
	)

	for i := 0; i < 10; i++ {
		ln, err = net.Listen("tcp", h.addr)
		if err == nil {
			h.logger.Info().Str("addr", h.addr).Msg("HTTP server listening")
			return ln
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	h.logger.Info().Msg("gracefully shutting down HTTP server")
	return h.Shutdown(ctx)
}

func (h *HTTPServer) handleSignals(ctx context.Context) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	defer h.Stop(ctx) // nolint:errcheck
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-signalCh:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				return
			default:
				panic(fmt.Sprintf("unsupported signal: %v", sig))
			}
		}
	}
}
