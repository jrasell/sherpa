package autoscale

import (
	"time"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	policyBackend "github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog"
)

type AutoScale struct {
	cfg    *Config
	logger zerolog.Logger
	nomad  *nomad.Client
	scaler scale.Scale

	policyBackend policyBackend.PolicyBackend
	pool          *ants.PoolWithFunc
	inProgress    bool

	// isRunning is used to track whether the autoscaler loop is being run. This helps determine
	// whether stop should be called.
	isRunning bool

	// doneChan is used to stop the autoscaling handler execution.
	doneChan chan struct{}
}

type scalableResources struct {
	cpu int
	mem int
}

type workerPayload struct {
	jobID  string
	policy map[string]*policy.GroupScalingPolicy
}

func NewAutoScaleServer(l zerolog.Logger, n *nomad.Client, p policyBackend.PolicyBackend, s scale.Scale, cfg *Config) (*AutoScale, error) {
	as := AutoScale{
		cfg:           cfg,
		logger:        l,
		nomad:         n,
		policyBackend: p,
		scaler:        s,
		doneChan:      make(chan struct{}),
	}

	pool, err := as.createWorkerPool()
	if err != nil {
		return nil, err
	}
	as.pool = pool

	return &as, nil
}

// IsRunning is used to determine if the autoscaler loop is running.
func (a *AutoScale) IsRunning() bool {
	return a.isRunning
}

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

				// Create a new policy object to track groups that are not considered to be in
				// deployment.
				nonDeploying := make(map[string]*policy.GroupScalingPolicy)

				// Iterate the group policies, and check whether they are in deployment or not.
				for group := range allPolicies[job] {
					if !a.scaler.JobGroupIsDeploying(job, group) {
						nonDeploying[group] = allPolicies[job][group]
					}
				}

				// If we have groups within the job that are not deploying, we can trigger a
				// scaling event.
				if len(nonDeploying) > 0 {
					if err := a.pool.Invoke(&workerPayload{jobID: job, policy: allPolicies[job]}); err != nil {
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
	return ants.NewPoolWithFunc(
		ants.Options{
			Capacity:       a.cfg.ScalingThreads,
			ExpiryDuration: 60 * time.Second,
			PoolFunc:       a.workerPoolFunc(),
		},
	)
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
		a.autoscaleJob(req.jobID, req.policy)
	}
}
