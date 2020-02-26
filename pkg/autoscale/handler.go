package autoscale

import (
	"time"

	"github.com/jrasell/sherpa/pkg/helper"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/autoscale/metrics"
	"github.com/jrasell/sherpa/pkg/autoscale/metrics/influxdb"
	"github.com/jrasell/sherpa/pkg/autoscale/metrics/prometheus"
	"github.com/jrasell/sherpa/pkg/policy"
	policyBackend "github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/scale"
	ants "github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog"
)

type AutoScale struct {
	cfg    *Config
	logger zerolog.Logger
	nomad  *nomad.Client
	scaler scale.Scale

	policyBackend policyBackend.PolicyBackend
	pool          *ants.PoolWithFunc

	// metricProvider
	metricProvider map[policy.MetricsProvider]metrics.Provider

	// isRunning is used to track whether the autoscaler loop is being run. This helps determine
	// whether stop should be called.
	isRunning bool

	// inProgress is used to determine if there is an autoscaling loop currently in progress.
	inProgress bool

	// doneChan is used to stop the autoscaling handler execution.
	doneChan chan struct{}
}

type workerPayload struct {
	time   time.Time
	jobID  string
	policy map[string]*policy.GroupScalingPolicy
}

func NewAutoScaleServer(cfg *SetupConfig) (*AutoScale, error) {
	as := AutoScale{
		cfg: &Config{
			ScalingInterval:   cfg.ScalingInterval,
			ScalingThreads:    cfg.ScalingThreads,
			StrictChecking:    cfg.StrictChecking,
			MetricProviderCfg: cfg.MetricProviderCfg,
		},
		logger:        cfg.Logger,
		nomad:         cfg.Nomad,
		policyBackend: cfg.PolicyBackend,
		scaler:        cfg.Scale,
		doneChan:      make(chan struct{}),
	}

	as.setupMetricProviders()

	pool, err := as.createWorkerPool()
	if err != nil {
		return nil, err
	}
	as.pool = pool

	return &as, nil
}

// setupMetricProviders setups up the metric providers.
func (a *AutoScale) setupMetricProviders() {

	// Initialise the metric provider map within AutoScale.
	a.metricProvider = make(map[policy.MetricsProvider]metrics.Provider)

	// If there is available Prometheus config, setup the provider.
	if a.cfg.MetricProviderCfg.Prometheus != nil {
		promClient, err := prometheus.NewClient(a.cfg.MetricProviderCfg.Prometheus.Addr, a.logger)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to setup Prometheus metric provider client")
		} else {
			a.metricProvider[policy.ProviderPrometheus] = promClient
		}
	}
	// If there is available InfluxDB config, setup the provider.
	if a.cfg.MetricProviderCfg.InfluxDB != nil {
		idbClient, err := influxdb.NewClient(
			a.cfg.MetricProviderCfg.InfluxDB.Addr,
			a.cfg.MetricProviderCfg.InfluxDB.Username,
			a.cfg.MetricProviderCfg.InfluxDB.Password,
			a.cfg.MetricProviderCfg.InfluxDB.Insecure,
			a.logger,
		)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to setup InfluxDB metric provider client")
		} else {
			a.metricProvider[policy.ProviderInfluxDB] = idbClient
		}
	}
}

// IsRunning is used to determine if the autoscaler loop is running.
func (a *AutoScale) IsRunning() bool {
	return a.isRunning
}

// Run starts the autoscaler ticker loop and only stops when Stop() is called.
func (a *AutoScale) Run() {
	a.logger.Info().Msg("starting Sherpa internal auto-scaling engine")

	// Track that the autoscaler is actively running.
	a.isRunning = true

	t := time.NewTicker(time.Second * time.Duration(a.cfg.ScalingInterval))
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// Check whether a previous scaling loop is in progress, and if it is we should skip
			// this round. This avoids putting more pressure on a system which may be under load
			// causing slow API responses.
			if a.inProgress {
				a.logger.Info().Msg("scaling run in progress, skipping new assessment")
				break
			}
			a.setScalingInProgressTrue()

			allPolicies, err := a.policyBackend.GetPolicies()
			if err != nil {
				a.logger.Error().Err(err).Msg("autoscaler unable to get scaling policies")
				a.setScalingInProgressFalse()
				break
			}
			totalPolicyCount := len(allPolicies)

			if totalPolicyCount == 0 {
				a.logger.Debug().Msg("no scaling policies found in storage backend")
				a.setScalingInProgressFalse()
				break
			}

			for job := range allPolicies {

				// Generate a timestamp for the occurrence of this autoscaling attempt.
				t := time.Now().UTC()

				// Create a new policy object to track groups that are not considered to be in
				// deployment or in cooldown.
				safeScale := make(map[string]*policy.GroupScalingPolicy)

				// Iterate the group policies, and check whether they are in deployment or in
				// cooldown.
				for group := range allPolicies[job] {

					// If the group policy is disabled, continue with the loop and ignore the
					// policy.
					if !allPolicies[job][group].Enabled {
						continue
					}

					// Deployment check.
					if a.scaler.JobGroupIsDeploying(job, group) {
						a.logger.Debug().
							Str("job", job).
							Str("group", group).
							Msg("job group is currently in deployment, skipping autoscaler evaluation")
						continue
					}

					// Cooldown check.
					cool, err := a.scaler.JobGroupIsInCooldown(job, group, allPolicies[job][group].Cooldown, t.UnixNano())
					if err != nil {
						a.logger.Error().
							Err(err).
							Str("job", job).
							Str("group", group).
							Msg("failed to determine if job group is in cooldown")
						continue
					}
					if cool {
						a.logger.Debug().
							Err(err).
							Str("job", job).
							Str("group", group).
							Msg("job group is currently in scaling cooldown, skipping autoscaler evaluation")
						continue
					}

					// At this point the initial checks have passed, therefore we can add the group
					// to the map indicating we can continue within the evaluation.
					safeScale[group] = allPolicies[job][group]
				}

				// If we have groups within the job that are not deploying, we can trigger a
				// scaling event.
				if len(safeScale) > 0 {
					if err := a.pool.Invoke(&workerPayload{jobID: job, policy: allPolicies[job], time: t}); err != nil {
						a.logger.Error().Err(err).Msg("failed to invoke autoscaling worker thread")
					}
				}
			}
			a.setScalingInProgressFalse()

		case <-a.doneChan:
			a.isRunning = false
			return
		}
	}
}

// Stop is used to gracefully stop the autoscaling workers.
func (a *AutoScale) Stop() {

	// Inform sub-process to exit.
	close(a.doneChan)

	for {
		if !a.isRunning && !a.inProgress {
			a.pool.Release()
			a.logger.Info().Msg("successfully drained autoscaler worker pool")
			return
		}
		a.logger.Debug().Msg("autoscaler still has in-flight workers, will continue to check")
		time.Sleep(1 * time.Second)
	}
}

func (a *AutoScale) setScalingInProgressTrue() {
	a.inProgress = true
}

func (a *AutoScale) setScalingInProgressFalse() {
	a.inProgress = false
}

// createWorkerPool is responsible for building the ants goroutine worker pool with the number of
// threads controlled by the operator configured value.
func (a *AutoScale) createWorkerPool() (*ants.PoolWithFunc, error) {
	return ants.NewPoolWithFunc(a.cfg.ScalingThreads, a.workerPoolFunc(), ants.WithExpiryDuration(60*time.Second))
}

func (a *AutoScale) workerPoolFunc() func(payload interface{}) {
	return func(payload interface{}) {

		// If this thread starts after the autoscaler has been asked to shutdown, exit. Otherwise
		// perform the work.
		select {
		case <-a.doneChan:
			a.logger.Debug().Msg("exiting autoscaling thread as a result of shutdown request")
			return
		default:
		}

		req, ok := payload.(*workerPayload)
		if !ok {
			a.logger.Error().Msg("autoscaler worker pool received unexpected payload type")
			return
		}

		newEval := autoscaleEvaluation{
			nomad:          a.nomad,
			metricProvider: a.metricProvider,
			scaler:         a.scaler,
			log:            helper.LoggerWithJobContext(a.logger, req.jobID),
			jobID:          req.jobID,
			policies:       req.policy,
			time:           req.time.UnixNano(),
		}
		newEval.evaluateJob()
	}
}
