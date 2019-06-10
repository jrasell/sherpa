package policy

// GroupScalingPolicy represents the configurable options for a task group scaling policy.
type GroupScalingPolicy struct {
	// Enabled is a boolean which tells whether the policy is disabled or not.
	Enabled bool `json:"Enabled"`

	// MinCount is the minimum count a task group should reach.
	MinCount int `json:"MinCount"`

	// MaxCount is the maximum count a task group should reach.
	MaxCount int `json:"MaxCount"`

	// ScaleOutCount is the number which a task group is incremented by during scaling.
	ScaleOutCount int `json:"ScaleOutCount"`

	// ScaleInCount is the number which a task group is decremented by during scaling.
	ScaleInCount int `json:"ScaleInCount"`

	ScaleOutCPUPercentageThreshold    int `json:"ScaleOutCPUPercentageThreshold"`
	ScaleOutMemoryPercentageThreshold int `json:"ScaleOutMemoryPercentageThreshold"`

	ScaleInCPUPercentageThreshold    int `json:"ScaleInCPUPercentageThreshold"`
	ScaleInMemoryPercentageThreshold int `json:"ScaleInMemoryPercentageThreshold"`
}

const (
	DefaultMinCount                          = 2
	DefaultMaxCount                          = 10
	DefaultScaleOutCount                     = 1
	DefaultScaleInCount                      = 1
	DefaultScaleOutCPUPercentageThreshold    = 80
	DefaultScaleOutMemoryPercentageThreshold = 80
	DefaultScaleInCPUPercentageThreshold     = 20
	DefaultScaleInMemoryPercentageThreshold  = 20
)
