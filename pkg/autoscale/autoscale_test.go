package autoscale

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/stretchr/testify/assert"
)

func Test_autoscaleEvaluation_buildScalingReq(t *testing.T) {
	testCases := []struct {
		ae             *autoscaleEvaluation
		inputDec       map[string]*scalingDecision
		expectedOutput []*scale.GroupReq
		name           string
	}{
		{
			ae: &autoscaleEvaluation{
				policies: map[string]*policy.GroupScalingPolicy{"test-group": {Cooldown: 3000}},
				time:     1313131313,
			},
			inputDec: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionIn,
					count:     13,
					metrics: map[string]*scalingMetricDecision{
						"datadog-cpu": {value: 99, threshold: 90},
						"nomad-cpu":   {value: 99, threshold: 90},
					},
				},
			},
			expectedOutput: []*scale.GroupReq{
				{
					Direction:          scale.DirectionIn,
					Count:              13,
					GroupName:          "test-group",
					GroupScalingPolicy: &policy.GroupScalingPolicy{Cooldown: 3000},
					Time:               1313131313,
					Meta: map[string]string{
						"datadog-cpu-threshold": "90.00",
						"datadog-cpu-value":     "99.00",
						"nomad-cpu-threshold":   "90.00",
						"nomad-cpu-value":       "99.00",
					},
				},
			},
			name: "single group decision",
		},

		{
			ae: &autoscaleEvaluation{
				policies: map[string]*policy.GroupScalingPolicy{
					"test-group": {Cooldown: 3000}, "test-group-1": {Cooldown: 3001},
				},
				time: 1313131313,
			},
			inputDec: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionIn,
					count:     13,
					metrics: map[string]*scalingMetricDecision{
						"datadog-cpu": {value: 99, threshold: 90},
						"nomad-cpu":   {value: 99, threshold: 90},
					},
				},
				"test-group-1": {
					direction: scale.DirectionOut,
					count:     131,
					metrics: map[string]*scalingMetricDecision{
						"prometheus-cpu": {value: 99, threshold: 90},
					},
				},
			},
			expectedOutput: []*scale.GroupReq{
				{
					Direction:          scale.DirectionIn,
					Count:              13,
					GroupName:          "test-group",
					GroupScalingPolicy: &policy.GroupScalingPolicy{Cooldown: 3000},
					Time:               1313131313,
					Meta: map[string]string{
						"datadog-cpu-threshold": "90.00",
						"datadog-cpu-value":     "99.00",
						"nomad-cpu-threshold":   "90.00",
						"nomad-cpu-value":       "99.00",
					},
				},
				{
					Direction:          scale.DirectionOut,
					Count:              131,
					GroupName:          "test-group-1",
					GroupScalingPolicy: &policy.GroupScalingPolicy{Cooldown: 3001},
					Time:               1313131313,
					Meta: map[string]string{
						"prometheus-cpu-threshold": "90.00",
						"prometheus-cpu-value":     "99.00",
					},
				},
			},
			name: "multiple group decision",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.ae.buildScalingReq(tc.inputDec)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_autoscaleEvaluation_buildSingleDecision(t *testing.T) {
	testCases := []struct {
		inputNomadDec  map[string]*scalingDecision
		inputExtDec    map[string]*scalingDecision
		expectedOutput map[string]*scalingDecision
		name           string
	}{
		{
			inputNomadDec: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionOut,
					count:     13,
					metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
				},
			},
			inputExtDec: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionOut,
					count:     13,
					metrics:   map[string]*scalingMetricDecision{"datadog-memory": {value: 99, threshold: 90}},
				},
			},
			expectedOutput: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionOut,
					count:     13,
					metrics: map[string]*scalingMetricDecision{
						"nomad-cpu":      {value: 99, threshold: 90},
						"datadog-memory": {value: 99, threshold: 90},
					},
				},
			},
			name: "single group, external and Nomad same direction",
		},

		{
			inputNomadDec: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionOut,
					count:     13,
					metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
				},
			},
			inputExtDec: map[string]*scalingDecision{
				"test-group-1": {
					direction: scale.DirectionIn,
					count:     13,
					metrics:   map[string]*scalingMetricDecision{"datadog-memory": {value: 99, threshold: 90}},
				},
			},
			expectedOutput: map[string]*scalingDecision{
				"test-group": {
					direction: scale.DirectionOut,
					count:     13,
					metrics: map[string]*scalingMetricDecision{
						"nomad-cpu": {value: 99, threshold: 90},
					},
				},
				"test-group-1": {
					direction: scale.DirectionIn,
					count:     13,
					metrics: map[string]*scalingMetricDecision{
						"datadog-memory": {value: 99, threshold: 90},
					},
				},
			},
			name: "external and Nomad different groups",
		},
	}
	ae := autoscaleEvaluation{}

	for _, tc := range testCases {
		actualOutput := ae.buildSingleDecision(tc.inputNomadDec, tc.inputExtDec)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_autoscaleEvaluation_buildSingleGroupDecision(t *testing.T) {
	testCases := []struct {
		inputNomadDec  *scalingDecision
		inputExtDec    *scalingDecision
		expectedOutput *scalingDecision
		name           string
	}{
		{
			inputNomadDec: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			inputExtDec: &scalingDecision{
				direction: scale.DirectionIn,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"datadog-cpu": {value: 99, threshold: 90}},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			name: "nomad decision out, external decision in",
		},
		{
			inputNomadDec: &scalingDecision{
				direction: scale.DirectionIn,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			inputExtDec: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"datadog-cpu": {value: 99, threshold: 90}},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"datadog-cpu": {value: 99, threshold: 90}},
			},
			name: "nomad decision in, external decision out",
		},
		{
			inputNomadDec: &scalingDecision{
				direction: scale.DirectionIn,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			inputExtDec: &scalingDecision{
				direction: scale.DirectionIn,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"datadog-cpu": {value: 99, threshold: 90}},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				count:     13,
				metrics: map[string]*scalingMetricDecision{
					"datadog-cpu": {value: 99, threshold: 90},
					"nomad-cpu":   {value: 99, threshold: 90},
				},
			},
			name: "nomad decision in, external decision in",
		},
		{
			inputNomadDec: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			inputExtDec: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"datadog-cpu": {value: 99, threshold: 90}},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics: map[string]*scalingMetricDecision{
					"datadog-cpu": {value: 99, threshold: 90},
					"nomad-cpu":   {value: 99, threshold: 90},
				},
			},
			name: "nomad decision out, external decision out",
		},
	}
	ae := autoscaleEvaluation{}

	for _, tc := range testCases {
		actualOutput := ae.buildSingleGroupDecision(tc.inputNomadDec, tc.inputExtDec)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_updateAutoscaleMeta(t *testing.T) {
	testCases := []struct {
		inputMetricType    string
		inputValue         float64
		inputThreshold     float64
		inputMeta          map[string]string
		expectedResultMeta map[string]string
		name               string
	}{
		{
			inputMetricType:    "nomad-cpu",
			inputValue:         90,
			inputThreshold:     89,
			inputMeta:          make(map[string]string),
			expectedResultMeta: map[string]string{"nomad-cpu-value": "90.00", "nomad-cpu-threshold": "89.00"},
			name:               "empty input meta",
		},
		{
			inputMetricType: "nomad-memory",
			inputValue:      91,
			inputThreshold:  89,
			inputMeta:       map[string]string{"nomad-cpu-value": "90.00", "nomad-cpu-threshold": "89.00"},
			expectedResultMeta: map[string]string{
				"nomad-cpu-value":        "90.00",
				"nomad-cpu-threshold":    "89.00",
				"nomad-memory-value":     "91.00",
				"nomad-memory-threshold": "89.00",
			},
			name: "existing input meta",
		},
	}

	for _, tc := range testCases {
		updateAutoscaleMeta(tc.inputMetricType, tc.inputValue, tc.inputThreshold, tc.inputMeta)
		assert.Equal(t, tc.expectedResultMeta, tc.inputMeta, tc.name)
	}
}
