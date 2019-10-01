package job

import (
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/watcher"
	"github.com/rs/zerolog"
)

type Watcher struct {
	logger          zerolog.Logger
	nomad           *api.Client
	lastChangeIndex uint64
}

func NewWatcher(logger zerolog.Logger, nomad *api.Client) watcher.Watcher {
	return &Watcher{
		logger: logger,
		nomad:  nomad,
	}
}

func (w *Watcher) Run(updateChan chan interface{}) {
	w.logger.Info().Msg("starting Sherpa Nomad meta policy watcher")

	var maxFound uint64

	q := &api.QueryOptions{WaitTime: 5 * time.Minute, WaitIndex: 1}

	for {

		jobs, meta, err := w.nomad.Jobs().List(q)
		if err != nil {
			w.logger.Error().Err(err).Msg("failed to call Nomad API for job listing")
			time.Sleep(10 * time.Second)
			continue
		}

		if !watcher.IndexHasChange(meta.LastIndex, q.WaitIndex) {
			w.logger.Debug().Msg("meta watcher last index has not changed")
			continue
		}
		w.logger.Debug().
			Uint64("old", q.WaitIndex).
			Uint64("new", meta.LastIndex).
			Msg("meta watcher last index has changed")

		// Iterate over all the returned jobs.
		for i := range jobs {

			// If the change index on the job is not newer than the previously recorded last index
			// we should continue to the next job. It is important here to use the lastChangeIndex
			// from the MetaWatcher as we want to process all jobs which have updated past this
			// index.
			if !watcher.IndexHasChange(jobs[i].ModifyIndex, w.lastChangeIndex) {
				continue
			}

			w.logger.Debug().
				Uint64("old", w.lastChangeIndex).
				Uint64("new", jobs[i].ModifyIndex).
				Str("job", jobs[i].ID).
				Msg("job modify index has changed is greater than last recorded")

			maxFound = watcher.MaxFound(jobs[i].ModifyIndex, maxFound)
			updateChan <- jobs[i]
		}

		// Update the Nomad API wait index to start long polling from the correct point and update
		// our recorded lastChangeIndex so we have the correct point to use during the next API
		// return.
		q.WaitIndex = meta.LastIndex
		w.lastChangeIndex = maxFound
	}
}
