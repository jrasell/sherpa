package scale

import (
	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
)

func (s *Scaler) sendScalingEventToState(job, id string, source state.Source, groupReqs []*GroupReq, err error) uuid.UUID {

	status := s.generateEventStatus(err)

	scaleID, err := uuid.NewV4()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to generate scaling UUID")
	}

	for i := range groupReqs {
		event := state.ScalingEventMessage{
			ID:        scaleID,
			EvalID:    id,
			GroupName: groupReqs[i].GroupName,
			Status:    status,
			Source:    source,
			Time:      groupReqs[i].Time,
			Count:     groupReqs[i].Count,
			Direction: groupReqs[i].Direction.String(),
		}

		if err := s.state.PutScalingEvent(job, &event); err != nil {
			s.logger.Error().
				Str("job", job).
				Str("group", event.GroupName).
				Err(err).Msg("failed to update state with scaling event")
		}
	}

	return scaleID
}

func (s *Scaler) generateEventStatus(err error) state.Status {
	switch err {
	case nil:
		return state.StatusCompleted
	default:
		return state.StatusFailed
	}
}
