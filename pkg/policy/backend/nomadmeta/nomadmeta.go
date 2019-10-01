package nomadmeta

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/policy/backend/memory"
	"github.com/rs/zerolog"
)

// NewJobScalingPolicies produces a new policy backend and processor. The policy backend is just
// the memory backend. The processor is used to handle job watcher updates, where the job is
// inspected for its status, and then any Sherpa meta parameters pulled out and validated.
func NewJobScalingPolicies(logger zerolog.Logger, nomad *api.Client) (backend.PolicyBackend, *Processor) {
	b := memory.NewJobScalingPolicies()
	return b, &Processor{
		logger:        logger,
		nomad:         nomad,
		backend:       b,
		jobUpdateChan: make(chan interface{}),
	}
}
