package autoscale

import (
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/pkg/errors"
)

type nomadGatheredMetrics struct {
	resourceInfo  map[string]*nomadResources
	resourceUsage map[string]*nomadResources
}

type nomadResources struct {
	cpu float64
	mem float64
}

const (
	nomadCPUMetricName    = "nomad-cpu"
	nomadMemoryMetricName = "nomad-memory"
)

// gatherNomadMetrics queries Nomad to produce Nomad resource allocation metrics for the job
// currently under evaluation. This only needs to be called once per job, and will provide stats
// for use across all groups.
func (ae *autoscaleEvaluation) gatherNomadMetrics() (*nomadGatheredMetrics, error) {
	resourceInfo, allocs, err := ae.getJobAllocations()
	if err != nil {
		return nil, err
	}

	// If there are no allocations which belong to a task group which has a scaling policy then we
	// can exit here. It is worth logging this detail as exiting here could be because a policy is
	// missing or has a typo within it.
	if len(allocs) < 1 {
		return nil, errors.New("no allocations found to match task group with scaling policy")
	}

	resourceUsage, err := ae.getJobResourceUsage(allocs)
	if err != nil {
		return nil, err
	}

	return &nomadGatheredMetrics{
		resourceInfo:  resourceInfo,
		resourceUsage: resourceUsage,
	}, nil
}

func (ae *autoscaleEvaluation) evaluateNomadJobMetrics(group string, pol *policy.GroupScalingPolicy, resources *nomadGatheredMetrics) *scalingDecision {

	// It is possible a scaling policy is configured for a job group, but the actual running
	// Nomad job doesn't have this job group configured. If this is the case, we should warn
	// the user in the logs and break the current loop.
	if _, ok := resources.resourceUsage[group]; !ok {
		ae.log.Warn().Str("group", group).Msg("job group found in policy but not found in Nomad job")
		return nil
	}

	// Maths. Find the current CPU and memory utilisation in percentage based on the total
	// available resources to the group, compared to their configured maximum based on the
	// resource stanza.
	cpuUsage := resources.resourceUsage[group].cpu * 100 / resources.resourceInfo[group].cpu
	memUsage := resources.resourceUsage[group].mem * 100 / resources.resourceInfo[group].mem
	ae.log.Info().
		Str("group", group).
		Float64("mem-value-percentage", memUsage).
		Float64("cpu-value-percentage", cpuUsage).
		Msg("Nomad resource utilisation calculation")

	return ae.calculateNomadScalingDecision(group, &nomadResources{cpu: cpuUsage, mem: memUsage}, pol)
}

func (ae *autoscaleEvaluation) getJobAllocations() (map[string]*nomadResources, []*nomad.Allocation, error) {
	out := make(map[string]*nomadResources)
	var allocList []*nomad.Allocation // nolint:prealloc

	allocs, _, err := ae.nomad.Jobs().Allocations(ae.jobID, false, nil)
	if err != nil {
		return out, nil, err
	}

	for i := range allocs {

		// GH-70: jobs can have a mix of groups with scaling policies, and groups without. We need
		// to safely check the policy.
		if v, ok := ae.policies[allocs[i].TaskGroup]; !ok {
			continue
		} else {
			if !v.Enabled {
				continue
			}
		}

		if !(allocs[i].ClientStatus == nomad.AllocClientStatusRunning || allocs[i].ClientStatus == nomad.AllocClientStatusPending) {
			continue
		}

		allocInfo, _, err := ae.nomad.Allocations().Info(allocs[i].ID, nil)
		if err != nil {
			return out, nil, err
		}
		allocList = append(allocList, allocInfo)

		updateResourceTracker(allocInfo.TaskGroup, float64(*allocInfo.Resources.CPU), float64(*allocInfo.Resources.MemoryMB), out)
	}
	return out, allocList, err
}

func (ae *autoscaleEvaluation) getJobResourceUsage(allocs []*nomad.Allocation) (map[string]*nomadResources, error) {
	out := make(map[string]*nomadResources)

	for i := range allocs {
		stats, err := ae.nomad.Allocations().Stats(allocs[i], nil)
		if err != nil {
			return out, err
		}

		updateResourceTracker(allocs[i].TaskGroup,
			stats.ResourceUsage.CpuStats.TotalTicks,
			float64(stats.ResourceUsage.MemoryStats.RSS/1024/1024),
			out)
	}
	return out, nil
}

// updateResourceTracker is responsible for updating the current resource tracking of a job, making
// sure nothing is overwritten where values already exists.
func updateResourceTracker(group string, cpu, mem float64, tracker map[string]*nomadResources) {
	if _, ok := tracker[group]; ok {
		tracker[group].mem += mem
		tracker[group].cpu += cpu
		return
	}
	tracker[group] = &nomadResources{cpu: cpu, mem: mem}
}
