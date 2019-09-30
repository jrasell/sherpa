package scale

import (
	"github.com/hashicorp/nomad/api"
)

// deploymentsKey is a composite key used for storing in-progress Nomad deployments.
type deploymentsKey struct {
	job, group string
}

// GetDeploymentChannel will return the channel where deployment updates should be sent.
func (s *Scaler) GetDeploymentChannel() chan interface{} {
	return s.deploymentUpdateChan
}

// JobGroupIsDeploying returns a boolean to indicate where or not the specified job and group is
// currently in deployment.
func (s *Scaler) JobGroupIsDeploying(job, group string) bool {
	s.deploymentsLock.RLock()
	_, ok := s.deployments[deploymentsKey{job: job, group: group}]
	s.deploymentsLock.RUnlock()
	return ok
}

// RunDeploymentUpdateHandler is used to handle updates and shutdowns when monitoring Nomad job
// deployments.
func (s *Scaler) RunDeploymentUpdateHandler() {
	s.logger.Info().Msg("starting scaler deployment update handler")

	for {
		select {
		case <-s.shutdownChan:
			return
		case msg := <-s.deploymentUpdateChan:
			go s.handleDeploymentMessage(msg)
		}
	}
}

func (s *Scaler) handleDeploymentMessage(msg interface{}) {
	deployment, ok := msg.(*api.Deployment)
	if !ok {
		s.logger.Error().Msg("received unexpected deployment update message type")
		return
	}
	s.logger.Debug().
		Str("status", deployment.Status).
		Str("job", deployment.JobID).
		Msg("received deployment update message to handle")

	s.deploymentsLock.Lock()
	defer s.deploymentsLock.Unlock()

	switch deployment.Status {
	case "running":
		// If the deployment is running, then we need to ensure that this is correctly tracked in
		// the scaler.
		for tg := range deployment.TaskGroups {
			s.deployments[deploymentsKey{job: deployment.JobID, group: tg}] = nil
		}

	default:
		// The default is used to catch paused, cancelled, failed, and successful deployments.
		// These result in the internal tracking of the deployment to be removed, indicating that
		// the job group is not in deployment and can therefore be scaled.
		for tg := range deployment.TaskGroups {
			delete(s.deployments, deploymentsKey{job: deployment.JobID, group: tg})
		}
	}
}
