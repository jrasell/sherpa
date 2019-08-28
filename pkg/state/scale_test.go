package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSource_String(t *testing.T) {
	testCases := []struct {
		source         Source
		expectedReturn string
		name           string
	}{
		{
			source:         SourceAPI,
			expectedReturn: "API",
			name:           "test API source",
		},
		{
			source:         SourceInternalAutoscaler,
			expectedReturn: "InternalAutoscaler",
			name:           "test InternalAutoscaler source",
		},
	}

	for _, tc := range testCases {
		actualReturn := tc.source.String()
		assert.Equal(t, tc.expectedReturn, actualReturn, tc.name)
	}
}

func TestStatus_String(t *testing.T) {
	testCases := []struct {
		status         Status
		expectedReturn string
		name           string
	}{
		{
			status:         StatusCompleted,
			expectedReturn: "Completed",
			name:           "test Completed status",
		},
		{
			status:         StatusFailed,
			expectedReturn: "Failed",
			name:           "test Failed status",
		},
	}

	for _, tc := range testCases {
		actualReturn := tc.status.String()
		assert.Equal(t, tc.expectedReturn, actualReturn, tc.name)
	}
}
