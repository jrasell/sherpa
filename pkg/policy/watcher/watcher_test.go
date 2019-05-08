package watcher

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMetaWatcher_policyFromMeta(t *testing.T) {
	watcher := NewMetaWatcher(zerolog.Logger{}, nil, nil)

	testCases := []struct {
		meta           map[string]string
		expectedPolicy *policy.GroupScalingPolicy
	}{
		{
			meta: map[string]string{
				metaKeyEnabled:       "true",
				metaKeyMaxCount:      "100",
				metaKeyMinCount:      "50",
				metaKeyScaleInCount:  "3",
				metaKeyScaleOutCount: "7",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:       true,
				MinCount:      50,
				MaxCount:      100,
				ScaleOutCount: 7,
				ScaleInCount:  3,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled: "true",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:       true,
				MinCount:      2,
				MaxCount:      10,
				ScaleOutCount: 1,
				ScaleInCount:  1,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled: "false",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:       false,
				MinCount:      2,
				MaxCount:      10,
				ScaleOutCount: 1,
				ScaleInCount:  1,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled:       "true",
				metaKeyMaxCount:      "10000",
				metaKeyScaleOutCount: "1000",
				metaKeyScaleInCount:  "10",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:       true,
				MinCount:      2,
				MaxCount:      10000,
				ScaleOutCount: 1000,
				ScaleInCount:  10,
			},
		},
		{
			meta: map[string]string{
				metaKeyEnabled:       "untranslatable",
				metaKeyMinCount:      "untranslatable",
				metaKeyMaxCount:      "untranslatable",
				metaKeyScaleOutCount: "untranslatable",
				metaKeyScaleInCount:  "untranslatable",
			},
			expectedPolicy: &policy.GroupScalingPolicy{
				Enabled:       false,
				MinCount:      2,
				MaxCount:      10,
				ScaleOutCount: 1,
				ScaleInCount:  1,
			},
		},
	}

	for _, tc := range testCases {
		policy := watcher.policyFromMeta(tc.meta)
		assert.Equal(t, tc.expectedPolicy, policy)
	}
}

func TestMetaWatcher_indexHasChange(t *testing.T) {
	watcher := NewMetaWatcher(zerolog.Logger{}, nil, nil)

	testCases := []struct {
		newValue       uint64
		oldValue       uint64
		expectedReturn bool
	}{
		{
			newValue:       13,
			oldValue:       7,
			expectedReturn: true,
		},
		{
			newValue:       13696,
			oldValue:       13696,
			expectedReturn: false,
		},
		{
			newValue:       7,
			oldValue:       13,
			expectedReturn: false,
		},
	}

	for _, tc := range testCases {
		res := watcher.indexHasChange(tc.newValue, tc.oldValue)
		assert.Equal(t, tc.expectedReturn, res)
	}
}
