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

	var maxFound uint64

	q := &api.QueryOptions{WaitTime: 5 * time.Minute, WaitIndex: 1}

	for {

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
		m.logger.Debug().
			Uint64("old", q.WaitIndex).
			Uint64("new", meta.LastIndex).
			Msg("meta watcher last index has changed")

		// Iterate over all the returned jobs.
		for i := range jobs {

			// If the change index on the job is not newer than the previously recorded last index
			// we should continue to the next job. It is important here to use the lastChangeIndex
			// from the MetaWatcher as we want to process all jobs which have updated past this
			// index.
			if !m.indexHasChange(jobs[i].ModifyIndex, m.lastChangeIndex) {
				continue
			}

			m.logger.Debug().
				Uint64("old", m.lastChangeIndex).
				Uint64("new", jobs[i].ModifyIndex).
				Str("job", jobs[i].ID).
				Msg("job modify index has changed is greater than last recorded")

			maxFound = m.maxFound(jobs[i].ModifyIndex, maxFound)
			go m.readJobMeta(jobs[i].ID)
		}

		// Update the Nomad API wait index to start long polling from the correct point and update
		// our recorded lastChangeIndex so we have the correct point to use during the next API
		// return.
		q.WaitIndex = meta.LastIndex
		m.lastChangeIndex = maxFound
	}
}

func (m *MetaWatcher) indexHasChange(new, old uint64) bool {
	if new <= old {
		return false
	}
	return true
}

func (m *MetaWatcher) maxFound(new, old uint64) uint64 {
	if new <= old {
		return old
	}
	return new
}
