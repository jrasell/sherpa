package nomadmeta

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/rs/zerolog"
)

type Processor struct {
	logger        zerolog.Logger
	nomad         *api.Client
	backend       backend.PolicyBackend
	jobUpdateChan chan interface{}
}

func (pr *Processor) Run() {
	pr.logger.Info().Msg("starting Nomad meta job update handler")
	for {
		select {
		case msg := <-pr.jobUpdateChan:
			go pr.handleJobListMessage(msg)
		}
	}
}

func (pr *Processor) GetUpdateChannel() chan interface{} {
	return pr.jobUpdateChan
}

func (pr *Processor) handleJobListMessage(msg interface{}) {
	job, ok := msg.(*api.JobListStub)
	if !ok {
		pr.logger.Error().Msg("received unexpected job update message type")
		return
	}
	pr.logger.Debug().Msg("received job list update message to handle")

	switch job.Status {
	case "running":
		go pr.handleRunningJob(job.ID)
	case "dead":
		go pr.handleDeadJob(job.ID)
	case "pending":
		// Pending is an in-between state, so just pass this through and do not do any work until
		// the job has a more actionable state.
	}
}

func (pr *Processor) handleDeadJob(jobID string) {
	if err := pr.backend.DeleteJobPolicy(jobID); err != nil {
		pr.logger.Error().
			Str("job", jobID).
			Err(err).
			Msg("failed to delete job group policies from backend store")
	}
}

func (pr *Processor) handleRunningJob(jobID string) {
	pr.logger.Debug().Str("job", jobID).Msg("reading job group meta stanzas")

	info, _, err := pr.nomad.Jobs().Info(jobID, nil)
	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to call Nomad API for job information")
		return
	}

	// Create a new object which will track all policies pulled from the job. Creating a new object
	// helps remove policies which have been removed from task groups as the policy state will be
	// overwritten.
	policies := map[string]*policy.GroupScalingPolicy{}

	for i := range info.TaskGroups {
		if pr.hasMetaKeys(info.TaskGroups[i].Meta) {
			policies[*info.TaskGroups[i].Name] = pr.policyFromMeta(info.TaskGroups[i].Meta)
		}
	}

	// If we have 0 policies, delete any stored policies for that job. This helps protect against
	// situations where a jobs meta scaling policy has been removed, but the job is still running.
	switch len(policies) {
	case 0:
		if err := pr.backend.DeleteJobPolicy(jobID); err != nil {
			pr.logger.Error().
				Str("job", jobID).
				Err(err).
				Msg("failed to delete job group policies from backend store")
		}
	default:
		if err := pr.backend.PutJobPolicy(jobID, policies); err != nil {
			pr.logger.Error().
				Str("job", jobID).
				Err(err).
				Msg("failed to add job group policies to backend store")
		}
	}
}

func (pr *Processor) policyFromMeta(meta map[string]string) *policy.GroupScalingPolicy {
	return &policy.GroupScalingPolicy{
		MaxCount:                          pr.maxCountValueOrDefault(meta),
		MinCount:                          pr.minCountValueOrDefault(meta),
		Enabled:                           pr.enabledValueOrDefault(meta),
		Cooldown:                          pr.cooldownValueOrDefault(meta),
		ScaleInCount:                      pr.scaleInValueOrDefault(meta),
		ScaleOutCount:                     pr.scaleOutValueOrDefault(meta),
		ScaleOutCPUPercentageThreshold:    pr.scaleOutCPUThresholdValueOrNil(meta),
		ScaleOutMemoryPercentageThreshold: pr.scaleOutMemoryThresholdValueOrNil(meta),
		ScaleInCPUPercentageThreshold:     pr.scaleInCPUThresholdValueOrNil(meta),
		ScaleInMemoryPercentageThreshold:  pr.scaleInMemoryThresholdValueOrNil(meta),
		ExternalChecks:                    pr.externalChecksFromMeta(meta),
	}
}

func (pr *Processor) enabledValueOrDefault(meta map[string]string) bool {
	if val, ok := meta[metaKeyEnabled]; ok {
		enabled, err := strconv.ParseBool(val)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert max count meta value to int")
			return false
		}
		return enabled
	}
	return false
}

func (pr *Processor) cooldownValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyCooldown]; ok {
		cooldown, err := strconv.Atoi(val)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert cooldown meta value to int")
			return policy.DefaultCooldown
		}
		return cooldown
	}
	return policy.DefaultCooldown
}

func (pr *Processor) maxCountValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyMaxCount]; ok {
		maxInt, err := strconv.Atoi(val)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert max count meta value to int")
			return policy.DefaultMaxCount
		}
		return maxInt
	}
	return policy.DefaultMaxCount
}

func (pr *Processor) minCountValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyMinCount]; ok {
		minInt, err := strconv.Atoi(val)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert min count meta value to int")
			return policy.DefaultMinCount
		}
		return minInt
	}
	return policy.DefaultMinCount
}

func (pr *Processor) scaleInValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleInCount]; ok {
		inCount, err := strconv.Atoi(val)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert scale in meta value to int")
			return policy.DefaultScaleInCount
		}
		return inCount
	}
	return policy.DefaultScaleInCount
}

func (pr *Processor) scaleOutValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleOutCount]; ok {
		outCount, err := strconv.Atoi(val)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert scale out meta value to int")
			return policy.DefaultScaleOutCount
		}
		return outCount
	}
	return policy.DefaultScaleOutCount
}

func (pr *Processor) scaleOutCPUThresholdValueOrNil(meta map[string]string) *float64 {
	if val, ok := meta[metaKeyScaleOutCPUPercentageThreshold]; ok {
		outThreshold, err := strconv.ParseFloat(val, 64)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert scale out CPU meta value to float64")
			return nil
		}
		return &outThreshold
	}
	return nil
}

func (pr *Processor) scaleInCPUThresholdValueOrNil(meta map[string]string) *float64 {
	if val, ok := meta[metaKeyScaleInCPUPercentageThreshold]; ok {
		outThreshold, err := strconv.ParseFloat(val, 64)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert scale in CPU meta value to float64")
			return nil
		}
		return &outThreshold
	}
	return nil
}

func (pr *Processor) scaleOutMemoryThresholdValueOrNil(meta map[string]string) *float64 {
	if val, ok := meta[metaKeyScaleOutMemoryPercentageThreshold]; ok {
		outThreshold, err := strconv.ParseFloat(val, 64)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert scale out memory meta value to float64")
			return nil
		}
		return &outThreshold
	}
	return nil
}

func (pr *Processor) scaleInMemoryThresholdValueOrNil(meta map[string]string) *float64 {
	if val, ok := meta[metaKeyScaleInMemoryPercentageThreshold]; ok {
		outThreshold, err := strconv.ParseFloat(val, 64)
		if err != nil {
			pr.logger.Error().Err(err).Msg("failed to convert scale in memory meta value to float64")
			return nil
		}
		return &outThreshold
	}
	return nil
}

func (pr *Processor) externalChecksFromMeta(meta map[string]string) map[string]*policy.ExternalCheck {
	if val, ok := meta[metaKeyExternalChecks]; ok {
		var checks map[string]*policy.ExternalCheck
		if err := json.Unmarshal([]byte(val), &checks); err != nil {
			pr.logger.Error().Err(err).Msg("failed to unmarshal external checks into struct")
			return nil
		}
		return checks
	}
	return nil
}

func (pr *Processor) hasMetaKeys(meta map[string]string) bool {
	if _, ok := meta[metaKeyEnabled]; ok {
		return true
	}
	return false
}
