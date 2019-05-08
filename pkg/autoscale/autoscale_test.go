package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoScale_updateResourceTracker(t *testing.T) {
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
