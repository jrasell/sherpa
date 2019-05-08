package v1

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/stretchr/testify/assert"
)

func Test_decodeGroupPolicyReqBodyAndValidate(t *testing.T) {
	testCases := []struct {
		body           []byte
		expectedPolicy *policy.GroupScalingPolicy
		expectedErr    error
	}{
		{
			body: []byte("{\"MaxCount\":10,\"MinCount\":2,\"Enabled\":true}"),
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:       true,
				MaxCount:      10,
				MinCount:      2,
				ScaleInCount:  1,
				ScaleOutCount: 1,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		policyRes, err := decodeGroupPolicyReqBodyAndValidate(tc.body)
		assert.Equal(t, tc.expectedPolicy, policyRes)

		if tc.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tc.expectedErr.Error())
		}
	}
}
