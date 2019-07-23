package autoscale

import (
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
)

func (a *AutoScale) autoscaleJob(jobID string, policies map[string]*policy.GroupScalingPolicy) {
	resourceInfo, allocs, err := a.getJobAllocations(jobID, policies)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("job", jobID).
			Msg("failed to gather allocation details for job")
		return
	}

	// If there are no allocations which belong to a task group which has a scaling policy then we
	// can exit here. It is worth logging this detail as exiting here could be because a policy is
	// missing or has a typo within it.
	if len(allocs) < 1 {
		a.logger.Debug().Str("job", jobID).Msg("no task groups in job have scaling policies enabled")
		return
	}

	resourceUsage, err := a.getJobResourceUsage(allocs)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("job", jobID).
			Msg("failed to gather job resource usage statistics")
		return
	}

	var scaleReq []*scale.GroupReq

	for group, pol := range policies {

		// It is possible a scaling policy is configured for a job group, but the actual running
		// Nomad job doesn't have this job group configured. If this is the case, we should warn
		// the user in the logs and break the current loop.
		if _, ok := resourceUsage[group]; !ok {
			a.logger.Warn().
				Str("job", jobID).
				Str("group", group).
				Msg("job group found in policy but not found in Nomad job")
			break
		}

		// Maths. Find the current CPU and memory utilisation in percentage based on the total
		// available resources to the group, compared to their configured maximum based on the
		// resource stanza.
		cpuUsage := resourceUsage[group].cpu / resourceInfo[group].cpu * 100
		memUsage := resourceUsage[group].mem / resourceInfo[group].mem * 100
		a.logger.Debug().
			Int("mem-usage-percentage", memUsage).
			Int("cpu-usage-percentage", cpuUsage).
			Msg("resource utilisation calculation")

		var scalingDir scale.Direction
		var count int

		switch {
		case cpuUsage < pol.ScaleInCPUPercentageThreshold, memUsage < pol.ScaleInMemoryPercentageThreshold:
			scalingDir = scale.DirectionIn
			count = pol.ScaleInCount
		case cpuUsage > pol.ScaleOutCPUPercentageThreshold, memUsage > pol.ScaleOutMemoryPercentageThreshold:
			scalingDir = scale.DirectionOut
			count = pol.ScaleOutCount
		}

		if scalingDir != "" {
			req := &scale.GroupReq{Direction: scalingDir, Count: count, GroupName: group, GroupScalingPolicy: pol}
			scaleReq = append(scaleReq, req)

			a.logger.Info().
				Str("job", jobID).
				Object("scaling-req", req).
				Msg("added group scaling request")
		}
	}

	// If group scaling requests have been added to the array for the job that is currently being
	// checked, trigger a scaling event.
	if len(scaleReq) > 0 {
		resp, _, err := a.scaler.Trigger(jobID, scaleReq)
		if err != nil {
			a.logger.Error().Str("job", jobID).Err(err).Msg("failed to trigger scaling of job")
		}

		if resp != nil {
			a.logger.Info().
				Str("job", jobID).
				Str("evaluation-id", resp.EvalID).
				Msg("successfully triggered autoscaling of job")
		}
	}
}

func (a *AutoScale) getJobAllocations(jobID string, policies map[string]*policy.GroupScalingPolicy) (map[string]*scalableResources, []*nomad.Allocation, error) {
	out := make(map[string]*scalableResources)
	var allocList []*nomad.Allocation // nolint:prealloc

	allocs, _, err := a.nomad.Jobs().Allocations(jobID, false, nil)
	for i := range allocs {

		if !policies[allocs[i].TaskGroup].Enabled {
			break
		}

		if !(allocs[i].ClientStatus == "running" || allocs[i].ClientStatus == "pending") {
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
