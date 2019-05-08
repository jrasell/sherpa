package watcher

import (
	"strconv"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/rs/zerolog"
)

type MetaWatcher struct {
	logger          zerolog.Logger
	nomad           *api.Client
	policies        backend.PolicyBackend
	lastChangeIndex uint64
}

func NewMetaWatcher(l zerolog.Logger, nomad *api.Client, p backend.PolicyBackend) *MetaWatcher {
	return &MetaWatcher{
		logger:   l,
		nomad:    nomad,
		policies: p,
	}
}

func (m *MetaWatcher) Run() {
	m.logger.Info().Msg("starting Sherpa Nomad meta policy engine")

	q := &api.QueryOptions{WaitTime: 5 * time.Minute}

	for {
		var maxFound uint64

		jobs, meta, err := m.nomad.Jobs().List(q)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to call Nomad API for job listing")
			time.Sleep(10 * time.Second)
			continue
		}

		if !m.indexHasChange(meta.LastIndex, q.WaitIndex) {
			m.logger.Debug().Msg("meta watcher last index has not changed")
			continue
		}

		m.logger.Debug().Msg("meta watcher last index has changed")
		maxFound = meta.LastIndex

		for i := range jobs {
			if !m.indexHasChange(jobs[i].ModifyIndex, maxFound) {
				continue
			}

			maxFound = jobs[i].ModifyIndex

			go m.readJobMeta(jobs[i].ID)
		}

		q.WaitIndex = maxFound
		m.lastChangeIndex = maxFound
	}
}

func (m *MetaWatcher) readJobMeta(jobID string) {
	info, _, err := m.nomad.Jobs().Info(jobID, nil)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to call Nomad API for job information")
		return
	}

	for i := range info.TaskGroups {
		if m.hasMetaKeys(info.TaskGroups[i].Meta) {
			p := m.policyFromMeta(info.TaskGroups[i].Meta)
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
		MaxCount:      m.maxCountValueOrDefault(meta),
		MinCount:      m.minCountValueOrDefault(meta),
		Enabled:       m.enabledValueOrDefault(meta),
		ScaleInCount:  m.scaleInValueOrDefault(meta),
		ScaleOutCount: m.scaleOutValueOrDefault(meta),
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

func (m *MetaWatcher) hasMetaKeys(meta map[string]string) bool {
	if _, ok := meta[metaKeyEnabled]; ok {
		return true
	}
	return false
}

func (m *MetaWatcher) indexHasChange(new, old uint64) bool {
	if new <= old {
		return false
	}
	return true
}
