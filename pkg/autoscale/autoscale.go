package autoscale

import (
	"fmt"
	"time"

	sendMetrics "github.com/armon/go-metrics"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/autoscale/metrics"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/rs/zerolog"
)

type autoscaleEvaluation struct {
	nomad          *nomad.Client
	metricProvider map[policy.MetricsProvider]metrics.Provider
	scaler         scale.Scale

	// policies are the job group policies that will be evaluated during this run.
	policies map[string]*policy.GroupScalingPolicy

	// jobID is the Nomad job which is under evaluation.
	jobID string

	// time is the unix nano time indicating when this scaling evaluation was triggered.
	time int64

	// log has the jobID context to save repeating this effort.
	log zerolog.Logger
}

func (ae *autoscaleEvaluation) evaluateJob() {
	ae.log.Debug().Msg("triggering autoscaling job evaluation")

	defer sendMetrics.MeasureSince([]string{"autoscale", ae.jobID, "evaluation"}, time.Now())

	externalDecision := make(map[string]*scalingDecision)
	nomadDecision := make(map[string]*scalingDecision)

	// We need to check to see whether the the job policies contain a group which is using Nomad
	// checks. This dictates whether we run the initial gatherNomadMetrics function and then
	// trigger the Nomad evaluation.
	var nomadCheck bool
	for _, p := range ae.policies {
		if p.NomadChecksEnabled() {
			nomadCheck = true
			break
		}
	}

	var (
		nomadMetricData *nomadGatheredMetrics
		err             error
	)

	// If the job policy contains groups which rely on Nomad data, we should collect this now. It
	// is most efficient to collect this data on a per job basis rather than per group. If we get
	// an error when performing this, log it and continue. It is possible external checks are also
	// in place and working; we can nil check the nomadMetricData to skip Nomad checks during this
	// evaluation.
	if nomadCheck {
		nomadMetricData, err = ae.gatherNomadMetrics()
		if err != nil {
			ae.log.Error().Err(err).Msg("failed to collect Nomad metrics, skipping Nomad based checks")
		}
	}

	// Iterate over the group policies for the job currently under evaluation.
	for group, p := range ae.policies {

		// Setup a start time so we can measure how long an individual job group evaluation takes.
		start := time.Now()
		ae.log.Debug().Str("group", group).Msg("triggering autoscaling job group evaluation")

		// If the group policy has Nomad checks enabled, and we managed to successfully get the
		// Nomad metric data, perform the evaluation.
		if p.NomadChecksEnabled() && nomadMetricData != nil {
			if nomadDec := ae.evaluateNomadJobMetrics(group, p, nomadMetricData); nomadDec != nil {
				nomadDecision[group] = nomadDec
			}
		}

		// If the group has external checks, perform these and ensure the decision if not nil,
		// before adding this to the decision tree.
		if p.ExternalChecks != nil {
			if extDec := ae.calculateExternalScalingDecision(group, p); extDec != nil {
				externalDecision[group] = extDec
			}
		}

		// This iteration has ended, so record the Sherpa metric.
		sendMetrics.MeasureSince([]string{"autoscale", ae.jobID, group, "evaluation"}, start)
	}

	ae.evaluateDecisions(nomadDecision, externalDecision)
}

func (ae *autoscaleEvaluation) evaluateDecisions(nomadDecision, externalDecision map[string]*scalingDecision) {

	// Exit quickly if there are now scaling decisions to process.
	if len(nomadDecision) == 0 && len(externalDecision) == 0 {
		ae.log.Info().Msg("scaling evaluation completed and no scaling required")
		return
	}
	var finalDecision map[string]*scalingDecision

	// Perform checks to see whether either the Nomad checks or the external checks have deemed
	// scaling to be required. A single decision here means we can easily move onto building the
	// scaling request.
	if len(nomadDecision) == 0 && len(externalDecision) > 0 {
		ae.log.Debug().Msg("scaling evaluation completed, handling scaling request based on external checks")
		finalDecision = externalDecision
	}
	if len(nomadDecision) > 0 && len(externalDecision) == 0 {
		ae.log.Debug().Msg("scaling evaluation completed, handling scaling request based on Nomad checks")
		finalDecision = nomadDecision
	}

	// If both Nomad and external sources believe the job needs scaling, we need to ensure that
	// they are not in different directions. If they are out will always trump in as we want to
	// ensure we can handle load.
	if len(nomadDecision) > 0 && len(externalDecision) > 0 {
		ae.log.Debug().
			Msg("scaling evaluation completed, handling scaling request based on Nomad and external checks")
		finalDecision = ae.buildSingleDecision(nomadDecision, externalDecision)
	}

	// Build the scaling request to send to the scaler backend.
	scaleReq := ae.buildScalingReq(finalDecision)

	// If group scaling requests have been added to the array for the job that is currently being
	// checked, trigger a scaling event. Run this in a routine as from this point there is nothing
	// we can do.
	if len(scaleReq) > 0 {
		go ae.triggerScaling(scaleReq)
	}
}

// triggerScaling is used to trigger the scaling of a job based on one or more group changes as
// as result of the scaling evaluation.
func (ae *autoscaleEvaluation) triggerScaling(req []*scale.GroupReq) {
	resp, _, err := ae.scaler.Trigger(ae.jobID, req, state.SourceInternalAutoscaler)
	if err != nil {
		ae.log.Error().Err(err).Msg("failed to trigger scaling of job")
		sendTriggerErrorMetrics(ae.jobID)
	}

	if resp != nil {
		ae.log.Info().
			Str("id", resp.ID.String()).
			Str("evaluation-id", resp.EvaluationID).
			Msg("successfully triggered autoscaling of job")
		sendTriggerSuccessMetrics(ae.jobID)
	}
}

// buildScalingReq takes the scaling decisions for the job under evaluation, and creates a list of
// group requests to send to the scaler backend.
func (ae *autoscaleEvaluation) buildScalingReq(dec map[string]*scalingDecision) []*scale.GroupReq {
	var scaleReq []*scale.GroupReq // nolint:prealloc

	for group, decision := range dec {

		// Iterate over the resource metrics which have broken their thresholds and ensure these
		// are added to the submission meta.
		meta := make(map[string]string)

		for name, metric := range decision.metrics {
			updateAutoscaleMeta(name, metric.value, metric.threshold, meta)
		}

		// Build the job group scaling request.
		req := &scale.GroupReq{
			Direction:          decision.direction,
			Count:              decision.count,
			GroupName:          group,
			GroupScalingPolicy: ae.policies[group],
			Time:               ae.time,
			Meta:               meta,
		}
		scaleReq = append(scaleReq, req)

		ae.log.Debug().
			Str("group", group).
			Object("scaling-req", req).
			Msg("added group scaling request")
	}
	return scaleReq
}

// buildSingleDecision takes the decisions from the Nomad and external providers checks, producing
// a single decision per job group.
func (ae *autoscaleEvaluation) buildSingleDecision(nomad, external map[string]*scalingDecision) map[string]*scalingDecision {
	final := make(map[string]*scalingDecision)

	for group, nomadDec := range nomad {
		if extDec, ok := external[group]; ok {
			final[group] = ae.buildSingleGroupDecision(nomadDec, extDec)
		}
		final[group] = nomadDec
	}

	for group, extDec := range external {
		if _, ok := final[group]; !ok {
			final[group] = extDec
		}
	}
	return final
}

// buildSingleGroupDecision takes a scalingDecision from Nomad and the external sources, deciding
// the final decision for the job group. In situations where out and in actions are requested, out
// will always win.
func (ae *autoscaleEvaluation) buildSingleGroupDecision(nomad, external *scalingDecision) *scalingDecision {

	// If Nomad wishes to scale in, but the external sources want to scale out, the external source
	// wins. We return this decision which contains the metric values which failed their threshold
	// check.
	if nomad.direction == scale.DirectionIn && external.direction == scale.DirectionOut {
		return external
	}

	// If Nomad wishes to scale out, but the external sources want to scale in, the Nomad source
	// wins. We return this decision which contains the metric values which failed their threshold
	// check.
	if nomad.direction == scale.DirectionOut && external.direction == scale.DirectionIn {
		return nomad
	}

	// If both Nomad and external sources desire the same action, combine the metrics which failed
	// their checks so we can provide this detail to the user.
	if nomad.direction == scale.DirectionIn && external.direction == scale.DirectionIn ||
		nomad.direction == scale.DirectionOut && external.direction == scale.DirectionOut {
		for key, metric := range external.metrics {
			nomad.metrics[key] = metric
		}
		return nomad
	}
	return nil
}

// updateAutoscaleMeta populates meta with the metrics used to autoscale a job group.
func updateAutoscaleMeta(metricType string, value, threshold float64, meta map[string]string) {
	key := metricType
	meta[key+"-value"] = fmt.Sprintf("%.2f", value)
	meta[key+"-threshold"] = fmt.Sprintf("%.2f", threshold)
}
