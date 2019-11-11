package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateResourceTracker(t *testing.T) {
	testCases := []struct {
		group    string
		cpu      float64
		mem      float64
		tracker  map[string]*nomadResources
		expected map[string]*nomadResources
	}{
		{
			group:    "cache",
			cpu:      100,
			mem:      100,
			tracker:  map[string]*nomadResources{"cache": {cpu: 100, mem: 100}},
			expected: map[string]*nomadResources{"cache": {cpu: 200, mem: 200}},
		},
		{
			group:    "cache",
			cpu:      200,
			mem:      200,
			tracker:  map[string]*nomadResources{},
			expected: map[string]*nomadResources{"cache": {cpu: 200, mem: 200}},
		},
	}

	for _, tc := range testCases {
		updateResourceTracker(tc.group, tc.cpu, tc.mem, tc.tracker)
		assert.Equal(t, tc.expected, tc.tracker)
	}
}
