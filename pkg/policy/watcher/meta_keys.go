package watcher

import (
	"strconv"

	"github.com/jrasell/sherpa/pkg/policy"
)

func (m *MetaWatcher) readJobMeta(jobID string) {
	m.logger.Debug().Str("job", jobID).Msg("reading job group meta stanzas")

	info, _, err := m.nomad.Jobs().Info(jobID, nil)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to call Nomad API for job information")
		return
	}

	for i := range info.TaskGroups {
		if m.hasMetaKeys(info.TaskGroups[i].Meta) {
			p := m.policyFromMeta(info.TaskGroups[i].Meta)

			// Launch a go-routine to attempt to write the job policy to the
			// backend store. If there is an error, we can't do anything but
			// log it allowing for investigation if needed.
			go func() {
				if err := m.policies.PutJobGroupPolicy(jobID, *info.TaskGroups[i].Name, p); err != nil {
					m.logger.Error().
						Str("job", jobID).
						Str("group", *info.TaskGroups[i].Name).
						Err(err).
						Msg("failed to add job group policy from Nomad meta")
				}
			}()
		}
	}
}

func (m *MetaWatcher) policyFromMeta(meta map[string]string) *policy.GroupScalingPolicy {
	return &policy.GroupScalingPolicy{
		MaxCount:                          m.maxCountValueOrDefault(meta),
		MinCount:                          m.minCountValueOrDefault(meta),
		Enabled:                           m.enabledValueOrDefault(meta),
		ScaleInCount:                      m.scaleInValueOrDefault(meta),
		ScaleOutCount:                     m.scaleOutValueOrDefault(meta),
		ScaleOutCPUPercentageThreshold:    m.scaleOutCPUThresholdValueOrDefault(meta),
		ScaleOutMemoryPercentageThreshold: m.scaleOutMemoryThresholdValueOrDefault(meta),
		ScaleInCPUPercentageThreshold:     m.scaleInCPUThresholdValueOrDefault(meta),
		ScaleInMemoryPercentageThreshold:  m.scaleInMemoryThresholdValueOrDefault(meta),
	}
}

func (m *MetaWatcher) enabledValueOrDefault(meta map[string]string) bool {
	if val, ok := meta[metaKeyEnabled]; ok {
		enabled, err := strconv.ParseBool(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert max count meta value to int")
			return false
		}
		return enabled
	}
	return false
}

func (m *MetaWatcher) maxCountValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyMaxCount]; ok {
		maxInt, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert max count meta value to int")
			return policy.DefaultMaxCount
		}
		return maxInt
	}
	return policy.DefaultMaxCount
}

func (m *MetaWatcher) minCountValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyMinCount]; ok {
		maxInt, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert min count meta value to int")
			return policy.DefaultMinCount
		}
		return maxInt
	}
	return policy.DefaultMinCount
}

func (m *MetaWatcher) scaleInValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleInCount]; ok {
		inCount, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert scale in meta value to int")
			return policy.DefaultScaleInCount
		}
		return inCount
	}
	return policy.DefaultScaleInCount
}

func (m *MetaWatcher) scaleOutValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleOutCount]; ok {
		outCount, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert scale out meta value to int")
			return policy.DefaultScaleOutCount
		}
		return outCount
	}
	return policy.DefaultScaleOutCount
}

func (m *MetaWatcher) scaleOutCPUThresholdValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleOutCPUPercentageThreshold]; ok {
		outThreshold, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert scale out CPU meta value to int")
			return policy.DefaultScaleOutCPUPercentageThreshold
		}
		return outThreshold
	}
	return policy.DefaultScaleOutCPUPercentageThreshold
}

func (m *MetaWatcher) scaleInCPUThresholdValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleInCPUPercentageThreshold]; ok {
		outThreshold, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert scale in CPU meta value to int")
			return policy.DefaultScaleInCPUPercentageThreshold
		}
		return outThreshold
	}
	return policy.DefaultScaleInCPUPercentageThreshold
}

func (m *MetaWatcher) scaleOutMemoryThresholdValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleOutMemoryPercentageThreshold]; ok {
		outThreshold, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert scale out memory meta value to int")
			return policy.DefaultScaleOutMemoryPercentageThreshold
		}
		return outThreshold
	}
	return policy.DefaultScaleOutMemoryPercentageThreshold
}

func (m *MetaWatcher) scaleInMemoryThresholdValueOrDefault(meta map[string]string) int {
	if val, ok := meta[metaKeyScaleInMemoryPercentageThreshold]; ok {
		outThreshold, err := strconv.Atoi(val)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to convert scale in memory meta value to int")
			return policy.DefaultScaleInMemoryPercentageThreshold
		}
		return outThreshold
	}
	return policy.DefaultScaleInMemoryPercentageThreshold
}

func (m *MetaWatcher) hasMetaKeys(meta map[string]string) bool {
	if _, ok := meta[metaKeyEnabled]; ok {
		return true
	}
	return false
}
