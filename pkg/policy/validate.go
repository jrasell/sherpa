package policy

import "github.com/pkg/errors"

// Validate is used to check that a specified policy is not just golang default init struct values.
func Validate(pol *GroupScalingPolicy) error {
	if pol.MinCount == 0 &&
		pol.MaxCount == 0 &&
		!pol.Enabled &&
		pol.ScaleInCount == 0 &&
		pol.ScaleOutCount == 0 {
		return errors.New("please specify non-default scaling policy")
	}
	return nil
}

// MergeWithDefaults merges a client specified scaling policy with Sherpa defaults to create a
// fully hydrated object.
func MergeWithDefaults(pol *GroupScalingPolicy) {
	if pol.MinCount == 0 {
		pol.MinCount = DefaultMinCount
	}
	if pol.MaxCount == 0 {
		pol.MaxCount = DefaultMaxCount
	}
	if pol.ScaleOutCount == 0 {
		pol.ScaleOutCount = DefaultScaleOutCount
	}
	if pol.ScaleInCount == 0 {
		pol.ScaleInCount = DefaultScaleInCount
	}
	if pol.ScaleOutCPUPercentageThreshold == 0 {
		pol.ScaleOutCPUPercentageThreshold = DefaultScaleOutCPUPercentageThreshold
	}
	if pol.ScaleOutMemoryPercentageThreshold == 0 {
		pol.ScaleOutMemoryPercentageThreshold = DefaultScaleOutMemoryPercentageThreshold
	}
	if pol.ScaleInCPUPercentageThreshold == 0 {
		pol.ScaleInCPUPercentageThreshold = DefaultScaleInCPUPercentageThreshold
	}
	if pol.ScaleInMemoryPercentageThreshold == 0 {
		pol.ScaleInMemoryPercentageThreshold = DefaultScaleInMemoryPercentageThreshold
	}
}
