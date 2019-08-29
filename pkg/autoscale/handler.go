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
	}

	pool, err := as.createWorkerPool()
	if err != nil {
		return nil, err
	}
	as.pool = pool

	return &as, nil
}

func (a *AutoScale) Run() {
	a.logger.Info().Msg("starting Sherpa internal auto-scaling engine")

	t := time.NewTicker(time.Second * time.Duration(a.cfg.ScalingInterval))
	defer t.Stop()

	defer a.pool.Release()

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
				if err := a.pool.Invoke(&workerPayload{jobID: job, policy: allPolicies[job]}); err != nil {
					a.logger.Error().Err(err).Msg("failed to invoke autoscaling worker thread")
				}
			}

			a.setScalingInProgressFalse()
		}
	}
}

func (a AutoScale) setScalingInProgressTrue() {
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
			PoolFunc: func(payload interface{}) {
				req, ok := payload.(*workerPayload)
				if !ok {
					a.logger.Error().Msg("autoscaler worker pool received unexpected payload type")
					return
				}
				a.autoscaleJob(req.jobID, req.policy)
			},
		},
	)
}
