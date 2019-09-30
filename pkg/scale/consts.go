package scale

import (
	"github.com/gofrs/uuid"
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/state"
)

// Scale is the interface used for scaling a Nomad job.
type Scale interface {
	// Trigger performs scaling of 1 or more job groups which belong to the same job.
	Trigger(string, []*GroupReq, state.Source) (*ScalingResponse, int, error)

	// GetDeploymentChannel is used to return the channel where updates to Nomad deployments should
	// be sent.
	GetDeploymentChannel() chan interface{}

	// RunDeploymentUpdateHandler is used to trigger the long running process which handles
	// messages sent to the deployment update channel.
	RunDeploymentUpdateHandler()

	// JobGroupIsDeploying checks internal references to identify if the queried job group is
	// currently in deployment.
	JobGroupIsDeploying(job, group string) bool

	checkJobGroupExists(*api.Job, string) *api.TaskGroup

	getNewGroupCount(*api.TaskGroup, *GroupReq) int
	checkNewGroupCount(int, *GroupReq) error
}

// GroupReq is a single item of scaling information for a single job group.
type GroupReq struct {

	// Direction is the scaling direction which the group is requested to change.
	Direction Direction

	// Count is the number by which to change the count by in the desired direction. This
	// information can also be found within the GroupScalingPolicy, but the API and the internal
	// autoscaler have different ways in which to process data so it is their responsibility to
	// populate this field for use and moves this logic away from the trigger.
	Count int

	// GroupName is the name of the job group to scale in this request.
	GroupName string

	// GroupScalingPolicy should include the job group scaling policy if it exists within the
	// Sherpa server.
	GroupScalingPolicy *policy.GroupScalingPolicy
}

type ScalingResponse struct {
	ID           uuid.UUID
	EvaluationID string
}
