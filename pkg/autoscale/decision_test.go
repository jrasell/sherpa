package autoscale

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_calculateScalingDecision(t *testing.T) {
	testCases := []struct {
		inputParams    *scalingCheckParams
		expectedOutput *scalingDecision
		name           string
	}{
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    50,
					memory: 50,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: nil,
			name:           "cpu scale none, memory scale none",
		},
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    80,
					memory: 50,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"cpu": {usage: 80, threshold: 75}},
			},
			name: "cpu scale out, memory scale none",
		},
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    50,
					memory: 80,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"memory": {usage: 80, threshold: 75}},
			},
			name: "cpu scale none, memory scale out",
		},
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    80,
					memory: 80,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"memory": {usage: 80, threshold: 75}, "cpu": {usage: 80, threshold: 75}},
			},
			name: "cpu scale out, memory scale out",
		},

		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    20,
					memory: 50,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"cpu": {usage: 20, threshold: 30}},
			},
			name: "cpu scale in, memory scale none",
		},
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    50,
					memory: 20,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"memory": {usage: 20, threshold: 30}},
			},
			name: "cpu scale none, memory scale in",
		},
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    20,
					memory: 20,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"memory": {usage: 20, threshold: 30}, "cpu": {usage: 20, threshold: 30}},
			},
			name: "cpu scale in, memory scale in",
		},
		{
			inputParams: &scalingCheckParams{
				resourceUsage: &scalingMetrics{
					cpu:    90,
					memory: 10,
				},
				logger: zerolog.Logger{},
				policy: &policy.GroupScalingPolicy{
					ScaleOutCount:                     2,
					ScaleInCount:                      2,
					ScaleOutCPUPercentageThreshold:    75,
					ScaleOutMemoryPercentageThreshold: 75,
					ScaleInCPUPercentageThreshold:     30,
					ScaleInMemoryPercentageThreshold:  30,
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     2,
				metrics:   map[string]*scalingMetricDecision{"cpu": {usage: 90, threshold: 75}},
			},
			name: "cpu scale out, memory scale in",
		},
	}

	for _, tc := range testCases {
		actualOutput := calculateScalingDecision(tc.inputParams)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}
