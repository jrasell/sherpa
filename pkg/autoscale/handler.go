package autoscale

import (
	"time"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/rs/zerolog"
)

type AutoScale struct {
	cfg    *Config
	logger zerolog.Logger
	nomad  *nomad.Client
	scaler scale.Scale

	policyBackend backend.PolicyBackend
	state         *state
}

type state struct {
	inProgress bool
}

type scalableResources struct {
	cpu int
	mem int
}

func NewAutoScaleServer(l zerolog.Logger, n *nomad.Client, p backend.PolicyBackend, cfg *Config) *AutoScale {
	return &AutoScale{
		cfg:           cfg,
		logger:        l,
		nomad:         n,
		policyBackend: p,
		scaler:        scale.NewScaler(n, l, cfg.StrictChecking),
		state:         &state{},
	}
}

func (a *AutoScale) Run() {
	a.logger.Info().Msg("starting Sherpa internal auto-scaling engine")

	t := time.NewTicker(time.Second * time.Duration(a.cfg.ScalingInterval))
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// Check whether a previous scaling loop is in progress, and if it is we should skip
			// this round. This avoids putting more pressure on a system which may be under load
			// causing slow API responses.
			if a.state.inProgress {
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
				go a.runScalingCycle(job, allPolicies[job])
			}
			a.setScalingInProgressFalse()
		}
	}
}

func (a AutoScale) setScalingInProgressTrue() {
	a.state.inProgress = true
}

func (a *AutoScale) setScalingInProgressFalse() {
	a.state.inProgress = false
}

func (a *AutoScale) runScalingCycle(job string, policy map[string]*policy.GroupScalingPolicy) {
	a.autoscaleJob(job, policy)
}
