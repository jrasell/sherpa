package v1

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_payloadOrPolicyCount(t *testing.T) {
	testCases := []struct {
		payloadCount        int
		policy              *policy.GroupScalingPolicy
		direction           scale.Direction
		expectedCountReturn int
		expectedErrorReturn error
	}{
		{
			payloadCount:        13,
			policy:              &policy.GroupScalingPolicy{},
			direction:           scale.DirectionIn,
			expectedCountReturn: 13,
			expectedErrorReturn: nil,
		},
		{
			payloadCount:        0,
			policy:              nil,
			direction:           scale.DirectionIn,
			expectedCountReturn: 0,
			expectedErrorReturn: errors.New("no policy configured, specify a count to scale by"),
		},
		{
			payloadCount:        0,
			policy:              &policy.GroupScalingPolicy{ScaleInCount: 3},
			direction:           scale.DirectionIn,
			expectedCountReturn: 3,
			expectedErrorReturn: nil,
		},
		{
			payloadCount:        0,
			policy:              &policy.GroupScalingPolicy{ScaleOutCount: 7},
			direction:           scale.DirectionOut,
			expectedCountReturn: 7,
			expectedErrorReturn: nil,
		},
	}

	for _, tc := range testCases {
		returnCount, err := payloadOrPolicyCount(tc.payloadCount, tc.policy, tc.direction)
		assert.Equal(t, tc.expectedCountReturn, returnCount)

		if tc.expectedErrorReturn != nil {
			assert.EqualError(t, err, tc.expectedErrorReturn.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}
