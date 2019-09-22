package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMember_verifyClusterName(t *testing.T) {
	testCases := []struct {
		operatorName  string
		stateName     string
		expectedError bool
		name          string
	}{
		{
			operatorName:  "sherpa-test-cluster",
			stateName:     "sherpa-test-cluster",
			expectedError: false,
			name:          "operator name matches state",
		},
		{
			operatorName:  "sherpa-test-cluster",
			stateName:     "sherpa-prod-cluster",
			expectedError: true,
			name:          "operator name does not match state",
		},
		{
			operatorName:  "",
			stateName:     "sherpa-test-cluster",
			expectedError: false,
			name:          "no operator passed name",
		},
	}

	for _, tc := range testCases {
		m := Member{clusterName: tc.operatorName}
		err := m.verifyClusterName(tc.stateName)

		if !tc.expectedError {
			assert.Nil(t, err, tc.name)
		} else {
			assert.Error(t, err, tc.name)
		}
	}
}
