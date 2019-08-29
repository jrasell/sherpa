package scale

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestScaler_getNewGroupCount(t *testing.T) {
	scaler := NewScaler(nil, zerolog.Logger{}, nil, false)

	testCases := []struct {
		taskGroup      *api.TaskGroup
		groupReq       *GroupReq
		expectedReturn int
	}{
		{
			taskGroup:      api.NewTaskGroup("cache", 7),
			groupReq:       &GroupReq{Direction: DirectionOut, Count: 3},
			expectedReturn: 10,
		},
		{
			taskGroup:      api.NewTaskGroup("cache", 4),
			groupReq:       &GroupReq{Direction: DirectionIn, Count: 2},
			expectedReturn: 2,
		},
	}

	for _, tc := range testCases {
		newCount := scaler.getNewGroupCount(tc.taskGroup, tc.groupReq)
		assert.Equal(t, tc.expectedReturn, newCount)
	}
}

func TestScaler_checkNewGroupCount(t *testing.T) {
	scaler := NewScaler(nil, zerolog.Logger{}, nil, true)

	testCases := []struct {
		newCount       int
		groupReq       *GroupReq
		expectedReturn error
	}{
		{
			newCount: 99,
			groupReq: &GroupReq{
				Direction: DirectionOut,
				GroupScalingPolicy: &policy.GroupScalingPolicy{
					MaxCount: 100,
				},
			},
			expectedReturn: nil,
		},
		{
			newCount: 101,
			groupReq: &GroupReq{
				Direction: DirectionOut,
				GroupScalingPolicy: &policy.GroupScalingPolicy{
					MaxCount: 100,
				},
			},
			expectedReturn: errors.New("scaling action will break job group maximum threshold"),
		},
		{
			newCount: 2,
			groupReq: &GroupReq{
				Direction: DirectionIn,
				GroupScalingPolicy: &policy.GroupScalingPolicy{
					MinCount: 2,
				},
			},
			expectedReturn: nil,
		},
		{
			newCount: 1,
			groupReq: &GroupReq{
				Direction: DirectionIn,
				GroupScalingPolicy: &policy.GroupScalingPolicy{
					MinCount: 2,
				},
			},
			expectedReturn: errors.New("scaling action will break job group minimum threshold"),
		},
	}

	for _, tc := range testCases {
		err := scaler.checkNewGroupCount(tc.newCount, tc.groupReq)
		if tc.expectedReturn != nil {
			assert.EqualError(t, err, tc.expectedReturn.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestScaler_jobGroupExists(t *testing.T) {
	scaler := NewScaler(nil, zerolog.Logger{}, nil, false)

	testCases := []struct {
		job            *api.Job
		group          string
		expectedReturn interface{}
	}{
		{
			job:            generateJobWithTargetGroup("sherpa-cache"),
			group:          "sherpa-cache",
			expectedReturn: api.NewTaskGroup("sherpa-cache", 1),
		},
		{
			job:            generateJobWithTargetGroup("sherpa-cache"),
			group:          "sherpa-db",
			expectedReturn: (*api.TaskGroup)(nil),
		},
	}

	for _, tc := range testCases {
		res := scaler.checkJobGroupExists(tc.job, tc.group)
		assert.Equal(t, tc.expectedReturn, res)
	}
}

func generateJobWithTargetGroup(groupName string) *api.Job {
	newJobName := "test"
	newJob := api.Job{ID: &newJobName}
	newJob.TaskGroups = append(newJob.TaskGroups, api.NewTaskGroup(groupName, 1))
	return &newJob
}
