package watcher

import (
	"time"

	"github.com/hashicorp/nomad/api"
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

func (m *MetaWatcher) indexHasChange(new, old uint64) bool {
	if new <= old {
		return false
	}
	return true
}
