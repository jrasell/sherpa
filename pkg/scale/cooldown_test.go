package scale

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/helper"
	"github.com/jrasell/sherpa/pkg/state"
	stateMemory "github.com/jrasell/sherpa/pkg/state/scale/memory"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestScaler_JobGroupIsInCooldown(t *testing.T) {
	testCases := []struct {
		inputJobName         string
		inputGroupName       string
		inputCoolDown        int
		inputTime            int64
		lastScalingEvent     *state.ScalingEventMessage
		expectedCooldownResp bool
		name                 string
	}{
		{
			inputJobName:         "test-job-1",
			inputGroupName:       "test-group-1",
			inputCoolDown:        180,
			inputTime:            helper.GenerateEventTimestamp(),
			expectedCooldownResp: true,
			name:                 "job group with policy cooldown set is in scaling cooldown",
			lastScalingEvent: &state.ScalingEventMessage{
				ID:        uuid.UUID{},
				GroupName: "test-group-1",
				EvalID:    "test",
				Source:    "test",
				Time:      helper.GenerateEventTimestamp(),
				Status:    "test",
				Count:     1,
				Direction: "in",
			},
		},
		{
			inputJobName:         "test-job-1",
			inputGroupName:       "test-group-1",
			inputTime:            helper.GenerateEventTimestamp(),
			expectedCooldownResp: false,
			name:                 "job group without previous scaling event",
			lastScalingEvent:     nil,
		},
		{
			inputJobName:         "test-job-1",
			inputGroupName:       "test-group-1",
			inputCoolDown:        300,
			inputTime:            helper.GenerateEventTimestamp(),
			expectedCooldownResp: false,
			name:                 "job group policy with cooldown but last event long ago",
			lastScalingEvent: &state.ScalingEventMessage{
				ID:        uuid.UUID{},
				GroupName: "test-group-1",
				EvalID:    "test",
				Source:    "test",
				Time:      helper.GenerateEventTimestamp() - 1000000000000,
				Status:    "test",
				Count:     1,
				Direction: "in",
			},
		},
	}

	for _, tc := range testCases {

		// Create a new Scaler for each test, to ensure to conflicting resources.
		sc := Scaler{logger: zerolog.Logger{}, nomadClient: nil, state: stateMemory.NewStateBackend(), strict: true}

		// Write the last event to check against if this isn't nil, meaning we do not have one.
		if tc.lastScalingEvent != nil {
			assert.Nil(t, sc.state.PutScalingEvent(tc.inputJobName, tc.lastScalingEvent), tc.name)
		}

		cooldown, err := sc.JobGroupIsInCooldown(tc.inputJobName, tc.inputGroupName, tc.inputCoolDown, tc.inputTime)
		assert.Nil(t, err, tc.name)
		assert.Equal(t, tc.expectedCooldownResp, cooldown, tc.name)
	}
}
