package scale

import (
	"net/http"
	"strings"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ Scale = (*Scaler)(nil)

// Scale is the interface used for scaling a Nomad job.
type Scale interface {
	// Trigger performs scaling of 1 or more job groups which belong to the same job.
	Trigger(string, []*GroupReq) (*api.JobRegisterResponse, int, error)

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

func (g *GroupReq) MarshalZerologObject(e *zerolog.Event) {
	e.Str("direction", g.Direction.String()).
		Int("count", g.Count).
		Str("group", g.GroupName)
}

type Scaler struct {
	logger      zerolog.Logger
	nomadClient *api.Client
	strict      bool
}

func NewScaler(c *api.Client, l zerolog.Logger, strictChecking bool) Scale {
	return &Scaler{
		logger:      l,
		nomadClient: c,
		strict:      strictChecking,
	}
}

// Trigger performs scaling of 1 or more job groups which belong to the same job.
//
// The return values indicate:
//		- the Nomad API job register response
//		- the HTTP return code, used for the Sherpa API
//		- any error
func (s *Scaler) Trigger(jobID string, groupReqs []*GroupReq) (*api.JobRegisterResponse, int, error) {

	// In order to submit a job for scaling we need to read the entire job back to Nomad as it does
	// not currently have convenience methods for changing job group counts.
	job, found, err := s.getJob(jobID)
	if !found && err == nil {
		s.logger.Info().Str("job", jobID).Msg("job not found to be running")
		return nil, http.StatusNotFound, errors.New("job not found")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var changes bool

	if s.strict {
		changes = s.triggerWithStrictChecking(job, groupReqs)
	} else {
		changes = s.triggerWithoutStrictChecking(job, groupReqs)
	}

	// If no job changes have been processed, inform the client as such. There are a number of
	// reasons this could happen which will be presented in the Sherpa logs if needed.
	if !changes {
		return nil, http.StatusNotModified, nil
	}

	resp, err := s.triggerNomadRegister(job)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}

func (s *Scaler) triggerWithStrictChecking(job *api.Job, groupReqs []*GroupReq) bool {
	var changes bool

	for i := range groupReqs {
		// During strict checking it is easiest and quickest to check whether the job group has an
		// enabled policy and exit the current loop if it is not enabled.
		if !groupReqs[i].GroupScalingPolicy.Enabled {
			s.logger.Debug().
				Str("job", *job.ID).
				Str("group", groupReqs[i].GroupName).
				Msg("job group scaling policy is disabled")
			break
		}

		// Check that the running job on the Nomad cluster has the job group and grab this. The
		// func loops so we might as well grab the group out here and save more loops.
		tg := s.checkJobGroupExists(job, groupReqs[i].GroupName)
		if tg == nil {
			break
		}

		// Important: when dealing with a Nomad job we are dealing with a pointer. In strict
		// checking we should check the count outside of the job before modifying the job as its
		// possible some task groups pass checks and have updates and some don't. In this situation
		// we still want to submit the updated job.
		newCount := s.getNewGroupCount(tg, groupReqs[i])
		if err := s.checkNewGroupCount(newCount, groupReqs[i]); err != nil {
			s.logger.Debug().
				Str("job", *job.ID).
				Str("group", groupReqs[i].GroupName).
				Msg(err.Error())
			break
		}

		// Once the check is completed, update the job group count and ensure changes are marked as
		// true.
		*tg.Count = newCount
		changes = true
	}

	return changes
}

func (s *Scaler) triggerWithoutStrictChecking(job *api.Job, groupReqs []*GroupReq) bool {
	var changes bool

	for i := range groupReqs {
		tg := s.checkJobGroupExists(job, groupReqs[i].GroupName)
		if tg == nil {
			break
		}

		// Once we have confirmed the job group exists within the running Nomad job, we can assume
		// there are changes to the job to submit to Nomad.
		changes = true

		// As we do not have strict checking, we can blindly update the task group count.
		*tg.Count = s.getNewGroupCount(tg, groupReqs[i])
	}

	return changes
}

func (s *Scaler) getNewGroupCount(taskGroup *api.TaskGroup, req *GroupReq) int {
	switch req.Direction {
	case DirectionIn:
		return *taskGroup.Count - req.Count
	case DirectionOut:
		return *taskGroup.Count + req.Count
	}
	return 0
}

func (s *Scaler) checkNewGroupCount(newCount int, req *GroupReq) error {
	switch req.Direction {
	case DirectionIn:
		if newCount < req.GroupScalingPolicy.MinCount {
			return errors.New("scaling action will break job group minimum threshold")
		}
	case DirectionOut:
		if newCount > req.GroupScalingPolicy.MaxCount {
			return errors.New("scaling action will break job group maximum threshold")
		}
	}
	return nil
}

// triggerNomadRegister is used to submit the updated job to the Nomad API.
func (s *Scaler) triggerNomadRegister(job *api.Job) (*api.JobRegisterResponse, error) {
	resp, _, err := s.nomadClient.Jobs().Register(job, nil)
	return resp, err
}

func (s *Scaler) getJob(jobID string) (*api.Job, bool, error) {
	job, _, err := s.nomadClient.Jobs().Info(jobID, nil)

	// If the job is not running on the cluster, the Nomad API will return an error which contains
	// the 404 not found message. We want to be able to tell the difference between a 404 and an
	// actual error calling the API so the check the error string returned.
	if err != nil && strings.Contains(err.Error(), "404") {
		s.logger.Info().Err(err).Msg("failed to find job requested for scaling within Nomad")
		return nil, false, nil
	}

	// If the error does not contain 404, we can assume this was an actual error in calling the
	// Nomad API.
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to call the Nomad jobs API")
		return nil, false, err
	}
	return job, true, nil
}

// checkJobGroupExists checks that the passed group is configured within the job spec that is passed to
// the function also. This is helpful when wanting to perform safety checks and comparisons
// between configured scaling policies and actual Nomad server job configuration.
func (s *Scaler) checkJobGroupExists(job *api.Job, group string) *api.TaskGroup {
	for i := range job.TaskGroups {
		if *job.TaskGroups[i].Name == group {
			return job.TaskGroups[i]
		}
	}
	s.logger.Warn().
		Str("job", *job.ID).
		Str("group", group).
		Msg("task group not found within running Nomad job")
	return nil
}
