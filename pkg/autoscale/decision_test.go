package autoscale

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/helper"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/stretchr/testify/assert"
)

func Test_autoscaleEvaluation_calculateNomadScalingDecision(t *testing.T) {
	testCases := []struct {
		inputPolicy    *policy.GroupScalingPolicy
		inputResource  *nomadResources
		inputGroup     string
		expectedOutput *scalingDecision
		name           string
	}{
		{
			inputPolicy: &policy.GroupScalingPolicy{
				ScaleOutCPUPercentageThreshold: helper.Float64ToPointer(90.00),
				ScaleInCPUPercentageThreshold:  helper.Float64ToPointer(89.98),
			},
			inputResource:  &nomadResources{cpu: 89.99},
			inputGroup:     "test-group",
			expectedOutput: nil,
			name:           "CPU Nomad checks only scaling not required",
		},
		{
			inputPolicy: &policy.GroupScalingPolicy{
				ScaleOutCPUPercentageThreshold:    helper.Float64ToPointer(90.00),
				ScaleInCPUPercentageThreshold:     helper.Float64ToPointer(89.98),
				ScaleOutMemoryPercentageThreshold: helper.Float64ToPointer(90.00),
				ScaleInMemoryPercentageThreshold:  helper.Float64ToPointer(89.98),
			},
			inputResource:  &nomadResources{cpu: 89.99, mem: 89.99},
			inputGroup:     "test-group",
			expectedOutput: nil,
			name:           "all Nomad checks scaling not required",
		},
	}

	for _, tc := range testCases {
		ae := autoscaleEvaluation{policies: map[string]*policy.GroupScalingPolicy{"test-group": tc.inputPolicy}}
		actualOutput := ae.calculateNomadScalingDecision(tc.inputGroup, tc.inputResource, tc.inputPolicy)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_autoscaleEvaluation_choseCorrectDecision(t *testing.T) {
	testCases := []struct {
		inputGroup     string
		inputDecisions map[scale.Direction]*scalingDecision
		expectedOutput *scalingDecision
		name           string
	}{
		{
			inputGroup: "test-group",
			inputDecisions: map[scale.Direction]*scalingDecision{
				scale.DirectionOut: {
					direction: scale.DirectionOut,
					count:     13,
					metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
				},
				scale.DirectionIn: {
					direction: scale.DirectionIn,
					metrics:   map[string]*scalingMetricDecision{"nomad-memory": {value: 90, threshold: 99}},
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			name: "both scale in and scale out decisions",
		},
		{
			inputGroup: "test-group",
			inputDecisions: map[scale.Direction]*scalingDecision{
				scale.DirectionOut: {
					direction: scale.DirectionOut,
					count:     13,
					metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			name: "scale out decision only",
		},
		{
			inputGroup: "test-group",
			inputDecisions: map[scale.Direction]*scalingDecision{
				scale.DirectionIn: {
					direction: scale.DirectionIn,
					metrics:   map[string]*scalingMetricDecision{"nomad-memory": {value: 90, threshold: 99}},
				},
			},
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				count:     13,
				metrics:   map[string]*scalingMetricDecision{"nomad-memory": {value: 90, threshold: 99}},
			},
			name: "scale in decision only",
		},
	}

	ae := autoscaleEvaluation{
		policies: map[string]*policy.GroupScalingPolicy{"test-group": {ScaleInCount: 13, ScaleOutCount: 13}},
	}

	for _, tc := range testCases {
		actualOutput := ae.choseCorrectDecision(tc.inputGroup, tc.inputDecisions)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_updateDecisionMap(t *testing.T) {
	testCases := []struct {
		inputNew       *scalingDecision
		inputCur       map[scale.Direction]*scalingDecision
		expectedResult map[scale.Direction]*scalingDecision
		inputName      string
		name           string
	}{
		{
			inputNew: &scalingDecision{
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			},
			inputCur: make(map[scale.Direction]*scalingDecision),
			expectedResult: map[scale.Direction]*scalingDecision{scale.DirectionOut: {
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			}},
			inputName: "nomad-cpu",
			name:      "empty current map",
		},
		{
			inputNew: &scalingDecision{
				direction: scale.DirectionIn,
				metrics:   map[string]*scalingMetricDecision{"nomad-memory": {value: 90, threshold: 99}},
			},
			inputCur: map[scale.Direction]*scalingDecision{scale.DirectionOut: {
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			}},
			expectedResult: map[scale.Direction]*scalingDecision{
				scale.DirectionOut: {
					direction: scale.DirectionOut,
					metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
				},
				scale.DirectionIn: {
					direction: scale.DirectionIn,
					metrics:   map[string]*scalingMetricDecision{"nomad-memory": {value: 90, threshold: 99}},
				},
			},
			inputName: "nomad-memory",
			name:      "populated current map, different direction input",
		},
		{
			inputNew: &scalingDecision{
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"nomad-memory": {value: 90, threshold: 99}},
			},
			inputCur: map[scale.Direction]*scalingDecision{scale.DirectionOut: {
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"nomad-cpu": {value: 99, threshold: 90}},
			}},
			expectedResult: map[scale.Direction]*scalingDecision{
				scale.DirectionOut: {
					direction: scale.DirectionOut,
					metrics: map[string]*scalingMetricDecision{
						"nomad-cpu":    {value: 99, threshold: 90},
						"nomad-memory": {value: 90, threshold: 99},
					},
				},
			},
			inputName: "nomad-memory",
			name:      "populated current map, same direction input",
		},
	}

	for _, tc := range testCases {
		updateDecisionMap(tc.inputNew, tc.inputName, tc.inputCur)
		assert.Equal(t, tc.expectedResult, tc.inputCur, tc.name)
	}
}

func Test_performGreaterThanCheck(t *testing.T) {
	testCases := []struct {
		inputValue     float64
		inputCheck     float64
		inputName      string
		inputAction    policy.ComparisonAction
		expectedOutput *scalingDecision
		name           string
	}{
		{
			inputValue:  90,
			inputCheck:  91,
			inputName:   "test_cpu_check",
			inputAction: policy.ActionScaleOut,
			expectedOutput: &scalingDecision{
				direction: scale.DirectionNone,
				metrics:   make(map[string]*scalingMetricDecision),
			},
			name: "expected scale none result",
		},
		{
			inputValue:  90.001,
			inputCheck:  90,
			inputName:   "test_check",
			inputAction: policy.ActionScaleOut,
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"test_check": {value: 90.001, threshold: 90}},
			},
			name: "expected scale out result with metric decision",
		},
		{
			inputValue:  111001.1,
			inputCheck:  111001.01,
			inputName:   "test_check",
			inputAction: policy.ActionScaleIn,
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				metrics:   map[string]*scalingMetricDecision{"test_check": {value: 111001.1, threshold: 111001.01}},
			},
			name: "expected scale in result with metric decision",
		},
	}

	for _, tc := range testCases {
		actualOutput := performGreaterThanCheck(tc.inputValue, tc.inputCheck, tc.inputName, tc.inputAction)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_performLessThanCheck(t *testing.T) {
	testCases := []struct {
		inputValue     float64
		inputCheck     float64
		inputName      string
		inputAction    policy.ComparisonAction
		expectedOutput *scalingDecision
		name           string
	}{
		{
			inputValue:  91,
			inputCheck:  90,
			inputName:   "test_cpu_check",
			inputAction: policy.ActionScaleOut,
			expectedOutput: &scalingDecision{
				direction: scale.DirectionNone,
				metrics:   make(map[string]*scalingMetricDecision),
			},
			name: "expected scale none result",
		},
		{
			inputValue:  90,
			inputCheck:  90.001,
			inputName:   "test_check",
			inputAction: policy.ActionScaleOut,
			expectedOutput: &scalingDecision{
				direction: scale.DirectionOut,
				metrics:   map[string]*scalingMetricDecision{"test_check": {value: 90, threshold: 90.001}},
			},
			name: "expected scale out result with metric decision",
		},
		{
			inputValue:  111001.01,
			inputCheck:  111001.1,
			inputName:   "test_check",
			inputAction: policy.ActionScaleIn,
			expectedOutput: &scalingDecision{
				direction: scale.DirectionIn,
				metrics:   map[string]*scalingMetricDecision{"test_check": {value: 111001.01, threshold: 111001.1}},
			},
			name: "expected scale in result with metric decision",
		},
	}

	for _, tc := range testCases {
		actualOutput := performLessThanCheck(tc.inputValue, tc.inputCheck, tc.inputName, tc.inputAction)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}
