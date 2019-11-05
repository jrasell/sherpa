package autoscale

import (
	"strconv"
	"time"

	"github.com/armon/go-metrics"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/helper"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/jrasell/sherpa/pkg/state"
)

func (a *AutoScale) autoscaleJob(jobID string, policies map[string]*policy.GroupScalingPolicy, t int64) {
	defer metrics.MeasureSince([]string{"autoscale", jobID, "evaluation"}, time.Now())

	// Create a new logger with the job in the context.
	jobLogger := helper.LoggerWithJobContext(a.logger, jobID)

	resourceInfo, allocs, err := a.getJobAllocations(jobID, policies)
	if err != nil {
		jobLogger.Error().Err(err).Msg("failed to gather allocation details for job")
		sendEvaluationErrorMetrics(jobID)
		return
	}

	// If there are no allocations which belong to a task group which has a scaling policy then we
	// can exit here. It is worth logging this detail as exiting here could be because a policy is
	// missing or has a typo within it.
	if len(allocs) < 1 {
		jobLogger.Debug().Msg("no task groups in job have scaling policies enabled")
		return
	}

	resourceUsage, err := a.getJobResourceUsage(allocs)
	if err != nil {
		jobLogger.Error().Err(err).Msg("failed to gather job resource usage statistics")
		sendEvaluationErrorMetrics(jobID)
		return
	}

	var scaleReq []*scale.GroupReq // nolint:prealloc

	for group, pol := range policies {

		// Create a new logger with the group in the context from the job logger.
		groupLogger := helper.LoggerWithGroupContext(jobLogger, group)

		// It is possible a scaling policy is configured for a job group, but the actual running
		// Nomad job doesn't have this job group configured. If this is the case, we should warn
		// the user in the logs and break the current loop.
		if _, ok := resourceUsage[group]; !ok {
			groupLogger.Warn().Msg("job group found in policy but not found in Nomad job")
			continue
		}

		// Maths. Find the current CPU and memory utilisation in percentage based on the total
		// available resources to the group, compared to their configured maximum based on the
		// resource stanza.
		cpuUsage := resourceUsage[group].cpu * 100 / resourceInfo[group].cpu
		memUsage := resourceUsage[group].mem * 100 / resourceInfo[group].mem
		groupLogger.Debug().
			Int("mem-usage-percentage", memUsage).
			Int("cpu-usage-percentage", cpuUsage).
			Msg("resource utilisation calculation")

		decision := calculateScalingDecision(&scalingCheckParams{
			resourceUsage: &scalingMetrics{cpu: cpuUsage, memory: memUsage},
			logger:        groupLogger,
			policy:        pol,
		})

		// Exit if no scaling is required.
		if decision == nil {
			groupLogger.Debug().Msg("no scaling required")
			continue
		}
		groupLogger.Info().
			Object("decision", decision).
			Msg("scaling decision made and action required")

		// Iterate over the resource metrics which have broken their thresholds and ensure these
		// are added to the submission meta.
		meta := make(map[string]string)

		for name, metric := range decision.metrics {
			updateAutoscaleMeta(group, name, metric.usage, metric.threshold, meta)
		}

		// Build the job group scaling request.
		req := &scale.GroupReq{
			Direction:          decision.direction,
			Count:              decision.count,
			GroupName:          group,
			GroupScalingPolicy: pol,
			Time:               t,
		}
		scaleReq = append(scaleReq, req)

		groupLogger.Debug().
			Object("scaling-req", req).
			Msg("added group scaling request")
	}

	// If group scaling requests have been added to the array for the job that is currently being
	// checked, trigger a scaling event.
	if len(scaleReq) > 0 {
		resp, _, err := a.scaler.Trigger(jobID, scaleReq, state.SourceInternalAutoscaler)
		if err != nil {
			jobLogger.Error().Err(err).Msg("failed to trigger scaling of job")
			sendTriggerErrorMetrics(jobID)
		}

		if resp != nil {
			jobLogger.Info().
				Str("id", resp.ID.String()).
				Str("evaluation-id", resp.EvaluationID).
				Msg("successfully triggered autoscaling of job")
			sendTriggerSuccessMetrics(jobID)
		}
	}
}

func (a *AutoScale) getJobAllocations(jobID string, policies map[string]*policy.GroupScalingPolicy) (map[string]*scalableResources, []*nomad.Allocation, error) {
	out := make(map[string]*scalableResources)
	var allocList []*nomad.Allocation // nolint:prealloc

	allocs, _, err := a.nomad.Jobs().Allocations(jobID, false, nil)
	if err != nil {
		return out, nil, err
	}

	for i := range allocs {

		// GH-70: jobs can have a mix of groups with scaling policies, and groups without. We need
		// to safely check the policy.
		if v, ok := policies[allocs[i].TaskGroup]; !ok {
			break
		} else {
			if !v.Enabled {
				break
			}
		}

		if !(allocs[i].ClientStatus == nomad.AllocClientStatusRunning || allocs[i].ClientStatus == nomad.AllocClientStatusPending) {
			continue
		}

		allocInfo, _, err := a.nomad.Allocations().Info(allocs[i].ID, nil)
		if err != nil {
			return out, nil, err
		}
		allocList = append(allocList, allocInfo)

		updateResourceTracker(allocInfo.TaskGroup, *allocInfo.Resources.CPU, *allocInfo.Resources.MemoryMB, out)
	}
	return out, allocList, err
}

func (a *AutoScale) getJobResourceUsage(allocs []*nomad.Allocation) (map[string]*scalableResources, error) {
	out := make(map[string]*scalableResources)

	for i := range allocs {
		stats, err := a.nomad.Allocations().Stats(allocs[i], nil)
		if err != nil {
			return out, err
		}

		updateResourceTracker(allocs[i].TaskGroup,
			int(stats.ResourceUsage.CpuStats.TotalTicks),
			int(stats.ResourceUsage.MemoryStats.RSS/1024/1024),
			out)
	}
	return out, nil
}

// updateResourceTracker is responsible for updating the current resource tracking of a job, making
// sure nothing is overwritten where values already exists.
func updateResourceTracker(group string, cpu, mem int, tracker map[string]*scalableResources) {
	if _, ok := tracker[group]; ok {
		tracker[group].mem += mem
		tracker[group].cpu += cpu
		return
	}
	tracker[group] = &scalableResources{cpu: cpu, mem: mem}
}

// updateAutoscaleMeta populates meta with the metrics used to autoscale a job group.
func updateAutoscaleMeta(group, metricType string, value, threshold int, meta map[string]string) {
	key := group + "-" + metricType
	meta[key+"-value"] = strconv.Itoa(value)
	meta[key+"-threshold"] = strconv.Itoa(threshold)
}
