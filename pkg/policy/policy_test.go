package policy

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/helper"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGroupScalingPolicy_Validate(t *testing.T) {
	testCases := []struct {
		policy         GroupScalingPolicy
		expectedOutput error
		name           string
	}{
		{
			policy:         GroupScalingPolicy{},
			expectedOutput: errors.New("please specify non-default scaling policy"),
			name:           "invalid and empty scaling policy",
		},
		{
			policy: GroupScalingPolicy{
				Enabled:       true,
				Cooldown:      100,
				MinCount:      10,
				MaxCount:      1000,
				ScaleOutCount: 1,
				ScaleInCount:  1,
			},
			expectedOutput: nil,
			name:           "valid core params without external check",
		},
		{
			policy: GroupScalingPolicy{
				Enabled:       true,
				Cooldown:      100,
				MinCount:      10,
				MaxCount:      1000,
				ScaleOutCount: 1,
				ScaleInCount:  1,
				ExternalChecks: map[string]*ExternalCheck{"test_external_check": {
					Enabled:            true,
					Provider:           ProviderPrometheus,
					Query:              "what_do_you_get_when_you_multiply_six_by_nine",
					ComparisonOperator: ComparisonGreaterThan,
					ComparisonValue:    42,
					Action:             ActionScaleIn,
				}},
			},
			expectedOutput: nil,
			name:           "valid core params with external check",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.policy.Validate()
		if tc.expectedOutput == nil {
			assert.Nil(t, actualOutput, tc.name)
		} else {
			assert.EqualError(t, actualOutput, tc.expectedOutput.Error(), tc.name)
		}
	}
}

func TestGroupScalingPolicy_NomadChecksEnabled(t *testing.T) {
	testCases := []struct {
		policy         GroupScalingPolicy
		expectedOutput bool
		name           string
	}{
		{
			policy:         GroupScalingPolicy{},
			expectedOutput: false,
			name:           "no nomad checks enabled",
		},
		{
			policy: GroupScalingPolicy{
				ScaleOutCPUPercentageThreshold:    helper.Float64ToPointer(80),
				ScaleOutMemoryPercentageThreshold: helper.Float64ToPointer(80),
				ScaleInCPUPercentageThreshold:     helper.Float64ToPointer(80),
				ScaleInMemoryPercentageThreshold:  helper.Float64ToPointer(80),
			},
			expectedOutput: true,
			name:           "nomad checks enabled",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.policy.NomadChecksEnabled()
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func TestGroupScalingPolicy_MergeWithDefaults(t *testing.T) {
	testCases := []struct {
		inputPolicy    GroupScalingPolicy
		expectedOutput GroupScalingPolicy
		name           string
	}{
		{
			inputPolicy:    GroupScalingPolicy{Cooldown: 130},
			expectedOutput: GroupScalingPolicy{Cooldown: 130, MinCount: 2, MaxCount: 10, ScaleOutCount: 1, ScaleInCount: 1},
			name:           "cooldown set, all others default",
		},
		{
			inputPolicy:    GroupScalingPolicy{MinCount: 13},
			expectedOutput: GroupScalingPolicy{Cooldown: 180, MinCount: 13, MaxCount: 10, ScaleOutCount: 1, ScaleInCount: 1},
			name:           "min count set, all others default",
		},
		{
			inputPolicy:    GroupScalingPolicy{MaxCount: 13},
			expectedOutput: GroupScalingPolicy{Cooldown: 180, MinCount: 2, MaxCount: 13, ScaleOutCount: 1, ScaleInCount: 1},
			name:           "max count set, all others default",
		},
		{
			inputPolicy:    GroupScalingPolicy{ScaleOutCount: 13},
			expectedOutput: GroupScalingPolicy{Cooldown: 180, MinCount: 2, MaxCount: 10, ScaleOutCount: 13, ScaleInCount: 1},
			name:           "scale out count set, all others default",
		},
		{
			inputPolicy:    GroupScalingPolicy{ScaleInCount: 13},
			expectedOutput: GroupScalingPolicy{Cooldown: 180, MinCount: 2, MaxCount: 10, ScaleOutCount: 1, ScaleInCount: 13},
			name:           "scale in count set, all others default",
		},
	}
	for _, tc := range testCases {
		actualOutput := tc.inputPolicy.MergeWithDefaults()
		assert.Equal(t, tc.expectedOutput, *actualOutput, tc.name)
	}
}

func TestMetricsProvider_String(t *testing.T) {
	testCases := []struct {
		inputProvider  MetricsProvider
		expectedOutput string
	}{
		{inputProvider: ProviderPrometheus, expectedOutput: "prometheus"},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputProvider.String()
		assert.Equal(t, tc.expectedOutput, actualOutput)
	}
}

func TestMetricsProvider_Validate(t *testing.T) {
	const fakeProvider MetricsProvider = "fake-provider"

	testCases := []struct {
		inputOperator  MetricsProvider
		expectedOutput error
	}{
		{inputOperator: ProviderPrometheus, expectedOutput: nil},
		{inputOperator: fakeProvider, expectedOutput: errors.Errorf("Provider %s is not a valid option", fakeProvider.String())},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputOperator.Validate()
		if tc.expectedOutput == nil {
			assert.Nil(t, actualOutput)
		} else {
			assert.EqualError(t, actualOutput, tc.expectedOutput.Error())
		}
	}
}

func TestComparisonOperator_String(t *testing.T) {
	testCases := []struct {
		inputOperator  ComparisonOperator
		expectedOutput string
	}{
		{inputOperator: ComparisonGreaterThan, expectedOutput: "greater-than"},
		{inputOperator: ComparisonLessThan, expectedOutput: "less-than"},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputOperator.String()
		assert.Equal(t, tc.expectedOutput, actualOutput)
	}
}

func TestComparisonOperator_Validate(t *testing.T) {
	const fakeOperator ComparisonOperator = "fake-operator"

	testCases := []struct {
		inputOperator  ComparisonOperator
		expectedOutput error
	}{
		{inputOperator: ComparisonGreaterThan, expectedOutput: nil},
		{inputOperator: ComparisonLessThan, expectedOutput: nil},
		{inputOperator: fakeOperator, expectedOutput: errors.Errorf("ComparisonOperator %s is not a valid option", fakeOperator.String())},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputOperator.Validate()
		if tc.expectedOutput == nil {
			assert.Nil(t, actualOutput)
		} else {
			assert.EqualError(t, actualOutput, tc.expectedOutput.Error())
		}
	}
}

func TestComparisonAction_String(t *testing.T) {
	testCases := []struct {
		inputAction    ComparisonAction
		expectedOutput string
	}{
		{inputAction: ActionScaleIn, expectedOutput: "scale-in"},
		{inputAction: ActionScaleOut, expectedOutput: "scale-out"},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputAction.String()
		assert.Equal(t, tc.expectedOutput, actualOutput)
	}
}

func TestComparisonAction_Validate(t *testing.T) {
	const fakeAction ComparisonAction = "fake-action"

	testCases := []struct {
		inputAction    ComparisonAction
		expectedOutput error
	}{
		{inputAction: ActionScaleIn, expectedOutput: nil},
		{inputAction: ActionScaleOut, expectedOutput: nil},
		{inputAction: fakeAction, expectedOutput: errors.Errorf("Action %s is not a valid option", fakeAction.String())},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputAction.Validate()
		if tc.expectedOutput == nil {
			assert.Nil(t, actualOutput)
		} else {
			assert.EqualError(t, actualOutput, tc.expectedOutput.Error())
		}
	}
}
