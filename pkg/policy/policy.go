package policy

import (
	"github.com/pkg/errors"
)

// GroupScalingPolicy represents the configurable options for a task group scaling policy. Each
// policy configures high level parameters to control the overview scaling, and then specific
// checks which are used to evaluate whether the job requires scaling or not.
type GroupScalingPolicy struct {

	// Enabled is a boolean which tells whether the policy is disabled or not.
	Enabled bool `json:"Enabled"`

	// Cooldown is a time period in seconds. Once a scaling action has been triggered on the
	// desired group, another action will not be triggered until the cooldown period has
	// passed.
	Cooldown int `json:"Cooldown"`

	// MinCount is the minimum count a task group should reach.
	MinCount int `json:"MinCount"`

	// MaxCount is the maximum count a task group should reach.
	MaxCount int `json:"MaxCount"`

	// ScaleOutCount is the number which a task group is incremented by during scaling.
	ScaleOutCount int `json:"ScaleOutCount"`

	// ScaleInCount is the number which a task group is decremented by during scaling.
	ScaleInCount int `json:"ScaleInCount"`

	// ScaleOutCPUPercentageThreshold is used to perform an upper bound check on the CPU resource
	// consumption of a job group based on Nomad obtained metrics. This value can be nil indicating
	// this check should not be performed.
	ScaleOutCPUPercentageThreshold *float64 `json:"ScaleOutCPUPercentageThreshold,omitempty"`

	// ScaleOutMemoryPercentageThreshold is used to perform an upper bound check on the memory
	// resource consumption of a job group based on Nomad obtained metrics. This value can be nil
	// indicating this check should not be performed.
	ScaleOutMemoryPercentageThreshold *float64 `json:"ScaleOutMemoryPercentageThreshold,omitempty"`

	// ScaleInCPUPercentageThreshold is used to perform an lower bound check on the CPU resource
	// consumption of a job group based on Nomad obtained metrics. This value can be nil indicating
	//	// this check should not be performed.
	ScaleInCPUPercentageThreshold *float64 `json:"ScaleInCPUPercentageThreshold,omitempty"`

	// ScaleInMemoryPercentageThreshold is used to perform an lower bound check on the memory
	// resource consumption of a job group based on Nomad obtained metrics. This value can be nil
	// indicating this check should not be performed.
	ScaleInMemoryPercentageThreshold *float64 `json:"ScaleInMemoryPercentageThreshold,omitempty"`

	// ExternalChecks represents metrics which are gathered from external sources for analysis
	// during scaling evaluations. They are keyed by a user specified name which is a free form
	// string and does not have any requirements which impact the running on the check itself.
	ExternalChecks map[string]*ExternalCheck `json:"ExternalChecks,omitempty"`
}

// ExternalCheck is an individual check of a metric from an external source. The check contains all
// information required by the metric provider and autoscaler to perform its analysis and decision
// making.
type ExternalCheck struct {

	// Enabled is a boolean flag to identify whether this specific external check should be
	// actively run or not.
	Enabled bool `json:"Enabled"`

	// Provider is the external provider source for the query to run against.
	Provider MetricsProvider `json:"Provider"`

	// Query is the string representation of the query that will be run against the external
	// provider to obtain a single int value.
	Query string `json:"Query"`

	// ComparisonOperator
	ComparisonOperator ComparisonOperator `json:"ComparisonOperator"`

	// ComparisonValue is the float64 value that will be compared to the result of the metric
	// query.
	ComparisonValue float64 `json:"ComparisonValue"`

	// Action is the scaling action that should be taken if the queried metric fails the comparison
	// check.
	Action ComparisonAction `json:"Action"`
}

// Validate performs a number of checks on the GroupScalingPolicy to ensure it is valid for use.
func (gsp GroupScalingPolicy) Validate() error {

	// Check whether all the core policy parameters are at Go defaults. If this is the case return
	// an error.
	if gsp.MinCount == 0 && gsp.MaxCount == 0 &&
		gsp.Cooldown == 0 && !gsp.Enabled &&
		gsp.ScaleInCount == 0 && gsp.ScaleOutCount == 0 {
		return errors.New("please specify non-default scaling policy")
	}

	// Iterate over the external checks and validate the required components. The first error is
	// returned, rather than collecting.
	for name, check := range gsp.ExternalChecks {
		if err := check.Provider.Validate(); err != nil {
			return errors.Wrap(err, "failed to validate check"+name)
		}

		if err := check.ComparisonOperator.Validate(); err != nil {
			return errors.Wrap(err, "failed to validate check"+name)
		}

		if err := check.Action.Validate(); err != nil {
			return errors.Wrap(err, "failed to validate check"+name)
		}
	}

	return nil
}

// NomadChecksEnabled helps determine whether the group policy ins configured to run scaling checks
// based on Nomad resource metrics.
func (gsp GroupScalingPolicy) NomadChecksEnabled() bool {
	if gsp.ScaleInMemoryPercentageThreshold == nil && gsp.ScaleOutMemoryPercentageThreshold == nil &&
		gsp.ScaleInCPUPercentageThreshold == nil && gsp.ScaleOutCPUPercentageThreshold == nil {
		return false
	}
	if *gsp.ScaleInMemoryPercentageThreshold == 0 && *gsp.ScaleOutMemoryPercentageThreshold == 0 &&
		*gsp.ScaleInCPUPercentageThreshold == 0 && *gsp.ScaleOutCPUPercentageThreshold == 0 {
		return false
	}
	return true
}

// MergeWithDefaults iterates the GroupScalingPolicy core parameters, merging this with default
// params where the user has not set some.
func (gsp GroupScalingPolicy) MergeWithDefaults() *GroupScalingPolicy {
	n := gsp

	if n.MinCount == 0 {
		n.MinCount = DefaultMinCount
	}
	if n.MaxCount == 0 {
		n.MaxCount = DefaultMaxCount
	}
	if n.Cooldown == 0 {
		n.Cooldown = DefaultCooldown
	}
	if n.ScaleOutCount == 0 {
		n.ScaleOutCount = DefaultScaleOutCount
	}
	if n.ScaleInCount == 0 {
		n.ScaleInCount = DefaultScaleInCount
	}
	return &n
}

// MetricsProvider represents the backend providers which can supply metric values for autoscaling.
type MetricsProvider string

// String returns the string form of the MetricsProvider.
func (mp MetricsProvider) String() string { return string(mp) }

// Validate checks the MetricsProvider is a valid and that it can be handled within the autoscaler.
func (mp MetricsProvider) Validate() error {
	switch mp {
	case ProviderPrometheus:
		return nil
	case ProviderInfluxDB:
		return nil
	default:
		return errors.Errorf("Provider %s is not a valid option", mp.String())
	}
}

const (
	// ProviderPrometheus is the Prometheus metrics backend.
	ProviderPrometheus MetricsProvider = "prometheus"
	// ProviderInfluxDB is the InfluxDB metrics backend.
	ProviderInfluxDB MetricsProvider = "influxdb"
)

// ComparisonOperator is the operator used when evaluating a metric value against a threshold.
type ComparisonOperator string

// String returns the string form of the ComparisonOperator.
func (co ComparisonOperator) String() string { return string(co) }

// Validate checks the ComparisonOperator is a valid and that it can be handled within the
// autoscaler.
func (co ComparisonOperator) Validate() error {
	switch co {
	case ComparisonGreaterThan, ComparisonLessThan:
		return nil
	default:
		return errors.Errorf("ComparisonOperator %s is not a valid option", co.String())
	}
}

const (
	ComparisonGreaterThan ComparisonOperator = "greater-than"
	ComparisonLessThan    ComparisonOperator = "less-than"
)

// ComparisonAction is the action to take if the metric breaks the threshold.
type ComparisonAction string

// String returns the string form of the ComparisonAction.
func (ca ComparisonAction) String() string { return string(ca) }

// Validate checks the ComparisonAction is a valid and that it can be handled within the
// autoscaler.
func (ca ComparisonAction) Validate() error {
	switch ca {
	case ActionScaleIn, ActionScaleOut:
		return nil
	default:
		return errors.Errorf("Action %s is not a valid option", ca.String())
	}
}

const (
	// ActionScaleIn performs a scale in operation.
	ActionScaleIn ComparisonAction = "scale-in"

	// ActionScaleOut performs a scale out operation.
	ActionScaleOut ComparisonAction = "scale-out"
)

const (
	DefaultMinCount      = 2
	DefaultMaxCount      = 10
	DefaultCooldown      = 180
	DefaultScaleOutCount = 1
	DefaultScaleInCount  = 1
)
