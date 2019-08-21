package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoScale_createWorkerPool(t *testing.T) {
	testCases := []struct {
		autoscaleStruct *AutoScale
		expectedThreads int
		testName        string
	}{
		{
			autoscaleStruct: &AutoScale{cfg: &Config{ScalingThreads: 100}},
			expectedThreads: 100,
			testName:        "check worker pool number of concurrent threads",
		},
	}

	for _, tc := range testCases {
		pool, err := tc.autoscaleStruct.createWorkerPool()

		assert.Nil(t, err, tc.testName)
		assert.Equal(t, tc.expectedThreads, pool.Cap(), tc.testName)
	}
}
