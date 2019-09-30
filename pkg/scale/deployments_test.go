package scale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScaler_JobGroupIsDeploying(t *testing.T) {
	testCases := []struct {
		storedMap      map[deploymentsKey]interface{}
		inputJob       string
		inputGroup     string
		expectedResult bool
		name           string
	}{
		{
			storedMap:      generateStoredMap(),
			inputJob:       "test-job-1",
			inputGroup:     "test-group-1",
			expectedResult: true,
			name:           "job group within stored map",
		},
		{
			storedMap:      generateStoredMap(),
			inputJob:       "test-job-2",
			inputGroup:     "test-group-2",
			expectedResult: false,
			name:           "job group not within stored map",
		},
		{
			storedMap:      generateStoredMap(),
			inputJob:       "test-job-1",
			inputGroup:     "test-group-2",
			expectedResult: false,
			name:           "job exists but group not within stored map",
		},
	}

	for _, tc := range testCases {
		s := Scaler{deployments: tc.storedMap}
		actualResult := s.JobGroupIsDeploying(tc.inputJob, tc.inputGroup)
		assert.Equal(t, tc.expectedResult, actualResult, tc.name)
	}
}

func generateStoredMap() map[deploymentsKey]interface{} {
	m := make(map[deploymentsKey]interface{})
	m[deploymentsKey{job: "test-job-1", group: "test-group-1"}] = nil
	return m
}
