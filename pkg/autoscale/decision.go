package autoscale

import (
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/rs/zerolog"
)

type scalingDecision struct {
	direction scale.Direction
	count     int
	metrics   map[string]*scalingMetricDecision
}

// scalingMetricDecision describes the metric value and threshold which resulted in the decision to
// require scaling. These are ultimately used in meta submission on scaling trigger and so become
// available to operators.
type scalingMetricDecision struct {
	value     float64
	threshold float64
}

// MarshalZerologObject is used to marshal a scaling decision for logging with zerolog.
func (sd *scalingDecision) MarshalZerologObject(e *zerolog.Event) {
	e.Str("direction", sd.direction.String()).Int("count", sd.count)

	dict := zerolog.Dict()

	for metric, val := range sd.metrics {
		dict.Dict(metric, zerolog.Dict().
			Float64("threshold-percentage", val.threshold).
			Float64("value-percentage", val.value))
	}
	e.Dict("resources", dict)
}

// calculateNomadScalingDecision is used to figure out the scaling decision for the group based on
// configured Nomad metric checks.
func (ae *autoscaleEvaluation) calculateNomadScalingDecision(group string, use *nomadResources, pol *policy.GroupScalingPolicy) *scalingDecision {
	decisions := make(map[scale.Direction]*scalingDecision)

	// If the policy has a CPU scale out threshold, run this check.
	if pol.ScaleOutCPUPercentageThreshold != nil {
		cpuOutDec := performGreaterThanCheck(use.cpu, *pol.ScaleOutCPUPercentageThreshold,
			nomadCPUMetricName, policy.ActionScaleOut)
		updateDecisionMap(cpuOutDec, nomadCPUMetricName, decisions)
	}

	// If the policy has a CPU scale in threshold, run this check.
	if pol.ScaleInCPUPercentageThreshold != nil {
		cpuInDec := performLessThanCheck(use.cpu, *pol.ScaleInCPUPercentageThreshold,
			nomadCPUMetricName, policy.ActionScaleIn)
		updateDecisionMap(cpuInDec, nomadCPUMetricName, decisions)
	}

	// If the policy has a memory scale out threshold, run this check.
	if pol.ScaleOutMemoryPercentageThreshold != nil {
		memOutDec := performGreaterThanCheck(use.mem, *pol.ScaleOutMemoryPercentageThreshold,
			nomadMemoryMetricName, policy.ActionScaleOut)
		updateDecisionMap(memOutDec, nomadMemoryMetricName, decisions)
	}

	// If the policy has a memory scale in threshold, run this check.
	if pol.ScaleInMemoryPercentageThreshold != nil {
		memInDec := performLessThanCheck(use.mem, *pol.ScaleInMemoryPercentageThreshold,
			nomadMemoryMetricName, policy.ActionScaleIn)
		updateDecisionMap(memInDec, nomadMemoryMetricName, decisions)
	}

	return ae.choseCorrectDecision(group, decisions)
}

// calculateExternalScalingDecision is used to perform the scaling decision for the group based on
// configured external metric checks.
func (ae *autoscaleEvaluation) calculateExternalScalingDecision(group string, pol *policy.GroupScalingPolicy) *scalingDecision {
	decisions := make(map[scale.Direction]*scalingDecision)

	// Iterate each external check configured within the job group scaling policy.
	for name, check := range pol.ExternalChecks {

		// If the check is not enabled, skip this.
		if !check.Enabled {
			continue
		}

		if checkDecision := ae.evaluateExternalMetric(name, check); checkDecision != nil {
			updateDecisionMap(checkDecision, name, decisions)
		}
	}
	return ae.choseCorrectDecision(group, decisions)
}

// evaluateExternalMetric is used to trigger the evaluation on a named external check. The function
// handles getting the metric value, and comparing it against the configured policy check params.
func (ae *autoscaleEvaluation) evaluateExternalMetric(name string, check *policy.ExternalCheck) *scalingDecision {

	// Check that the provider is available and properly configured for use.
	if _, ok := ae.metricProvider[check.Provider]; !ok {
		ae.log.Warn().
			Str("metric-provider", check.Provider.String()).
			Msg("provider not found configured within autoscaler")
		return nil
	}

	// Perform the query to gather the metric value.
	value, err := ae.metricProvider[check.Provider].GetValue(check.Query)
	if err != nil {
		ae.log.Error().
			Err(err).
			Str("metric-provider", check.Provider.String()).
			Str("metric-query", check.Query).
			Msg("failed to query external provider for metric value")
		return nil
	}
	ae.log.Info().
		Err(err).
		Str("metric-provider", check.Provider.String()).
		Str("metric-query", check.Query).
		Float64("metric-value", *value).
		Msg("successfully queried external provider for metric value")

	switch check.ComparisonOperator {
	case policy.ComparisonGreaterThan:
		return performGreaterThanCheck(*value, check.ComparisonValue, name, check.Action)
	case policy.ComparisonLessThan:
		return performLessThanCheck(*value, check.ComparisonValue, name, check.Action)
	default:
		return nil
	}
}

// choseCorrectDecision takes a set of decisions made about the scaling direction of the group,
// and produces a single correct answer. This is mostly in place to ensure safety in situations
// where two different metric checks produce an out and an in decision.
func (ae *autoscaleEvaluation) choseCorrectDecision(group string, dec map[scale.Direction]*scalingDecision) *scalingDecision {

	// Always perform this check first to ensure out takes precedent over in.
	if dec[scale.DirectionIn] != nil && dec[scale.DirectionOut] != nil {
		ae.log.Info().Str("group", group).Msg("both scale in and scale out actions desired, using out action")
		dec[scale.DirectionOut].count = ae.policies[group].ScaleOutCount
		return dec[scale.DirectionOut]
	}

	if dec[scale.DirectionOut] != nil {
		dec[scale.DirectionOut].count = ae.policies[group].ScaleOutCount
		return dec[scale.DirectionOut]
	}

	if dec[scale.DirectionIn] != nil {
		dec[scale.DirectionIn].count = ae.policies[group].ScaleInCount
		return dec[scale.DirectionIn]
	}
	return nil
}

// updateDecisionMap is used to safely update a decision mapping based on the new decision.
func updateDecisionMap(new *scalingDecision, name string, cur map[scale.Direction]*scalingDecision) {
	if _, ok := cur[new.direction]; !ok {
		cur[new.direction] = &scalingDecision{
			direction: new.direction, metrics: make(map[string]*scalingMetricDecision),
		}
	}
	cur[new.direction].metrics[name] = new.metrics[name]
}

func performGreaterThanCheck(value, check float64, name string, action policy.ComparisonAction) *scalingDecision {
	resp := scalingDecision{metrics: make(map[string]*scalingMetricDecision), direction: scale.DirectionNone}

	if value > check {
		switch action {
		case policy.ActionScaleIn:
			resp.direction = scale.DirectionIn
			resp.metrics[name] = &scalingMetricDecision{value: value, threshold: check}
		case policy.ActionScaleOut:
			resp.direction = scale.DirectionOut
			resp.metrics[name] = &scalingMetricDecision{value: value, threshold: check}
		default:
		}
	}
	return &resp
}

func performLessThanCheck(value, check float64, name string, action policy.ComparisonAction) *scalingDecision {
	resp := scalingDecision{metrics: make(map[string]*scalingMetricDecision), direction: scale.DirectionNone}

	if value < check {
		switch action {
		case policy.ActionScaleIn:
			resp.direction = scale.DirectionIn
			resp.metrics[name] = &scalingMetricDecision{value: value, threshold: check}
		case policy.ActionScaleOut:
			resp.direction = scale.DirectionOut
			resp.metrics[name] = &scalingMetricDecision{value: value, threshold: check}
		default:
		}
	}
	return &resp
}
