package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateResourceTracker(t *testing.T) {
	testCases := []struct {
		group    string
		cpu      int
		mem      int
		tracker  map[string]*scalableResources
		expected map[string]*scalableResources
	}{
		{
			group:    "cache",
			cpu:      100,
			mem:      100,
			tracker:  map[string]*scalableResources{"cache": {cpu: 100, mem: 100}},
			expected: map[string]*scalableResources{"cache": {cpu: 200, mem: 200}},
		},
		{
			group:    "cache",
			cpu:      200,
			mem:      200,
			tracker:  map[string]*scalableResources{},
			expected: map[string]*scalableResources{"cache": {cpu: 200, mem: 200}},
		},
	}

	for _, tc := range testCases {
		updateResourceTracker(tc.group, tc.cpu, tc.mem, tc.tracker)
		assert.Equal(t, tc.expected, tc.tracker)
	}
}

func Test_updateAutoscaleMeta(t *testing.T) {
	testCases := []struct {
		inputGroup         string
		inputMetricType    string
		inputValue         int
		inputThreshold     int
		inputMeta          map[string]string
		expectedResultMeta map[string]string
		name               string
	}{
		{
			inputGroup:         "test",
			inputMetricType:    "cpu",
			inputValue:         90,
			inputThreshold:     89,
			inputMeta:          make(map[string]string),
			expectedResultMeta: map[string]string{"test-cpu-value": "90", "test-cpu-threshold": "89"},
			name:               "empty input meta",
		},
		{
			inputGroup:      "test",
			inputMetricType: "memory",
			inputValue:      91,
			inputThreshold:  89,
			inputMeta:       map[string]string{"test-cpu-value": "90", "test-cpu-threshold": "89"},
			expectedResultMeta: map[string]string{
				"test-cpu-value":        "90",
				"test-cpu-threshold":    "89",
				"test-memory-value":     "91",
				"test-memory-threshold": "89",
			},
			name: "existing input meta",
		},
	}

	for _, tc := range testCases {
		updateAutoscaleMeta(tc.inputGroup, tc.inputMetricType, tc.inputValue, tc.inputThreshold, tc.inputMeta)
		assert.Equal(t, tc.expectedResultMeta, tc.inputMeta, tc.name)
	}
}
