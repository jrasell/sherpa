package autoscale

import (
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/rs/zerolog"
)

// scalingCheckParams is used to perform a decision check on a particular job group.
type scalingCheckParams struct {
	resourceUsage *scalingMetrics
	logger        zerolog.Logger
	policy        *policy.GroupScalingPolicy
}

type scalingMetrics struct {
	cpu    int
	memory int
}

type scalingDecision struct {
	direction scale.Direction
	count     int
	metrics   map[string]*scalingMetricDecision
}

type scalingMetricDecision struct {
	usage     int
	threshold int
}

// calculateScalingDecision is used to determine whether or not a scaling action is desired for the
// job group in question.
func calculateScalingDecision(params *scalingCheckParams) *scalingDecision {
	var inDecision, outDecision *scalingDecision

	// Grab our decisions for analysing.
	inDecision = isScalingInRequired(params)
	outDecision = isScalingOutRequired(params)

	// Check if it has been decided we should scale both out and in. This happens when a groups
	// resource settings are suboptimal and one resource parameter is over-provisioned, and the
	// other is under-provisioned. Performing this check allows operators to clearly see this is
	// happening and potentially take action. It is possible in the future, this could be a feature
	// of Sherpa.
	if inDecision.direction != scale.DirectionNone && outDecision.direction != scale.DirectionNone {
		params.logger.Info().Msg("both scale in and scale out actions desired, using out action")
		return outDecision
	}

	if outDecision.direction != scale.DirectionNone {
		return outDecision
	}

	if inDecision.direction != scale.DirectionNone {
		return inDecision
	}

	return nil
}

func isScalingOutRequired(params *scalingCheckParams) *scalingDecision {
	resp := scalingDecision{metrics: make(map[string]*scalingMetricDecision), direction: scale.DirectionNone}

	// Perform a check to see if scaling in is required based on CPU utilisation.
	if params.resourceUsage.cpu > params.policy.ScaleOutCPUPercentageThreshold {
		resp.metrics["cpu"] = &scalingMetricDecision{
			usage:     params.resourceUsage.cpu,
			threshold: params.policy.ScaleOutCPUPercentageThreshold,
		}
		resp.direction = scale.DirectionOut
		resp.count = params.policy.ScaleOutCount
	}

	// Perform a check to see if scaling in is required based on memory utilisation.
	if params.resourceUsage.memory > params.policy.ScaleOutMemoryPercentageThreshold {
		resp.metrics["memory"] = &scalingMetricDecision{
			usage:     params.resourceUsage.memory,
			threshold: params.policy.ScaleOutMemoryPercentageThreshold,
		}
		resp.direction = scale.DirectionOut
		resp.count = params.policy.ScaleOutCount
	}

	return &resp
}

func isScalingInRequired(params *scalingCheckParams) *scalingDecision {
	resp := scalingDecision{metrics: make(map[string]*scalingMetricDecision), direction: scale.DirectionNone}

	// Perform a check to see if scaling in is required based on CPU utilisation.
	if params.resourceUsage.cpu < params.policy.ScaleInCPUPercentageThreshold {
		resp.metrics["cpu"] = &scalingMetricDecision{
			usage:     params.resourceUsage.cpu,
			threshold: params.policy.ScaleInCPUPercentageThreshold,
		}
		resp.direction = scale.DirectionIn
		resp.count = params.policy.ScaleInCount
	}

	// Perform a check to see if scaling in is required based on memory utilisation.
	if params.resourceUsage.memory < params.policy.ScaleInMemoryPercentageThreshold {
		resp.metrics["memory"] = &scalingMetricDecision{
			usage:     params.resourceUsage.memory,
			threshold: params.policy.ScaleInMemoryPercentageThreshold,
		}
		resp.direction = scale.DirectionIn
		resp.count = params.policy.ScaleInCount
	}

	return &resp
}

// MarshalZerologObject is used to marshal a scaling decision for logging with zerolog.
func (sd *scalingDecision) MarshalZerologObject(e *zerolog.Event) {
	e.Str("direction", sd.direction.String()).Int("count", sd.count)

	dict := zerolog.Dict()

	for metric, val := range sd.metrics {
		dict.Dict(metric, zerolog.Dict().
			Int("threshold-percentage", val.threshold).
			Int("usage-percentage", val.usage))
	}
	e.Dict("resources", dict)
}
