package nomadmeta

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestProcessor_policyFromMeta(t *testing.T) {
	_, p := NewJobScalingPolicies(zerolog.Logger{}, nil)

	testCases := []struct {
		meta           map[string]string
		expectedPolicy *policy.GroupScalingPolicy
	}{
		{
			meta: map[string]string{
				metaKeyEnabled:                           "true",
				metaKeyMaxCount:                          "100",
				metaKeyMinCount:                          "50",
				metaKeyScaleInCount:                      "3",
				metaKeyScaleOutCount:                     "7",
				metaKeyScaleOutCPUPercentageThreshold:    "95",
				metaKeyScaleOutMemoryPercentageThreshold: "95",
				metaKeyScaleInCPUPercentageThreshold:     "55",
				metaKeyScaleInMemoryPercentageThreshold:  "55",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:                           true,
				MinCount:                          50,
				MaxCount:                          100,
				ScaleOutCount:                     7,
				ScaleInCount:                      3,
				ScaleOutCPUPercentageThreshold:    95,
				ScaleOutMemoryPercentageThreshold: 95,
				ScaleInCPUPercentageThreshold:     55,
				ScaleInMemoryPercentageThreshold:  55,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled: "true",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:                           true,
				MinCount:                          2,
				MaxCount:                          10,
				ScaleOutCount:                     1,
				ScaleInCount:                      1,
				ScaleOutCPUPercentageThreshold:    80,
				ScaleOutMemoryPercentageThreshold: 80,
				ScaleInCPUPercentageThreshold:     20,
				ScaleInMemoryPercentageThreshold:  20,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled: "false",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:                           false,
				MinCount:                          2,
				MaxCount:                          10,
				ScaleOutCount:                     1,
				ScaleInCount:                      1,
				ScaleOutCPUPercentageThreshold:    80,
				ScaleOutMemoryPercentageThreshold: 80,
				ScaleInCPUPercentageThreshold:     20,
				ScaleInMemoryPercentageThreshold:  20,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled:                           "true",
				metaKeyMaxCount:                          "10000",
				metaKeyScaleOutCount:                     "1000",
				metaKeyScaleInCount:                      "10",
				metaKeyScaleOutCPUPercentageThreshold:    "95",
				metaKeyScaleOutMemoryPercentageThreshold: "75",
				metaKeyScaleInCPUPercentageThreshold:     "15",
				metaKeyScaleInMemoryPercentageThreshold:  "35",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:                           true,
				MinCount:                          2,
				MaxCount:                          10000,
				ScaleOutCount:                     1000,
				ScaleInCount:                      10,
				ScaleOutCPUPercentageThreshold:    95,
				ScaleOutMemoryPercentageThreshold: 75,
				ScaleInCPUPercentageThreshold:     15,
				ScaleInMemoryPercentageThreshold:  35,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled:                           "untranslatable",
				metaKeyMinCount:                          "untranslatable",
				metaKeyMaxCount:                          "untranslatable",
				metaKeyScaleOutCount:                     "untranslatable",
				metaKeyScaleInCount:                      "untranslatable",
				metaKeyScaleOutCPUPercentageThreshold:    "untranslatable",
				metaKeyScaleOutMemoryPercentageThreshold: "untranslatable",
				metaKeyScaleInCPUPercentageThreshold:     "untranslatable",
				metaKeyScaleInMemoryPercentageThreshold:  "untranslatable",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:                           false,
				MinCount:                          2,
				MaxCount:                          10,
				ScaleOutCount:                     1,
				ScaleInCount:                      1,
				ScaleOutCPUPercentageThreshold:    80,
				ScaleOutMemoryPercentageThreshold: 80,
				ScaleInCPUPercentageThreshold:     20,
				ScaleInMemoryPercentageThreshold:  20,
			},
		},
	}

	for _, tc := range testCases {
		actualPolicy := p.policyFromMeta(tc.meta)
		assert.Equal(t, tc.expectedPolicy, actualPolicy)
	}
}
