package policy

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	testCases := []struct {
		policy         *GroupScalingPolicy
		expectedReturn error
	}{
		{policy: &GroupScalingPolicy{}, expectedReturn: errors.New("please specify non-default scaling policy")},
		{policy: &GroupScalingPolicy{Enabled: true}, expectedReturn: nil},
	}

	for _, tc := range testCases {
		actualReturn := Validate(tc.policy)
		if tc.expectedReturn != nil {
			assert.EqualError(t, actualReturn, tc.expectedReturn.Error())
		} else {
			assert.Nil(t, actualReturn)
		}
	}
}

func Test_MergeWithDefaults(t *testing.T) {
	testCases := []struct {
		policy               *GroupScalingPolicy
		expectedResultPolicy *GroupScalingPolicy
	}{
		{
			policy: &GroupScalingPolicy{Enabled: true},
			expectedResultPolicy: &GroupScalingPolicy{
				Enabled:                           true,
				MinCount:                          DefaultMinCount,
				MaxCount:                          DefaultMaxCount,
				ScaleInCount:                      DefaultScaleInCount,
				ScaleOutCount:                     DefaultScaleOutCount,
				ScaleOutCPUPercentageThreshold:    DefaultScaleOutCPUPercentageThreshold,
				ScaleInCPUPercentageThreshold:     DefaultScaleInCPUPercentageThreshold,
				ScaleOutMemoryPercentageThreshold: DefaultScaleOutMemoryPercentageThreshold,
				ScaleInMemoryPercentageThreshold:  DefaultScaleInMemoryPercentageThreshold,
			},
		},
	}

	for _, tc := range testCases {
		MergeWithDefaults(tc.policy)
		assert.Equal(t, tc.expectedResultPolicy, tc.policy)
	}
}
