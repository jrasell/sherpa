package deployment

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

func New(logger zerolog.Logger, nomad *api.Client) watcher.Watcher {
	return &Watcher{
		logger: logger,
		nomad:  nomad,
	}
}

func (w *Watcher) Run(updateChan chan interface{}) {
	w.logger.Info().Msg("starting Sherpa Nomad deployment watcher")

	var maxFound uint64

	q := &api.QueryOptions{WaitTime: 5 * time.Minute, WaitIndex: 1}

	for {

		deployments, meta, err := w.nomad.Deployments().List(q)
		if err != nil {
			w.logger.Error().Err(err).Msg("failed to call Nomad API for deployment listing")
			time.Sleep(10 * time.Second)
			continue
		}

		if !watcher.IndexHasChange(meta.LastIndex, q.WaitIndex) {
			w.logger.Debug().Msg("deployment watcher last index has not changed")
			continue
		}
		w.logger.Debug().
			Uint64("old", q.WaitIndex).
			Uint64("new", meta.LastIndex).
			Msg("deployment watcher last index has changed")

		// Iterate over all the returned deployments.
		for i := range deployments {

			// If the change index on the deployment is not newer than the previously recorded last
			// index we should continue to the next deployment. It is important here to use the
			// lastChangeIndex from the watcher as we want to process all deployments which have
			// updated past this index.
			if !watcher.IndexHasChange(deployments[i].ModifyIndex, w.lastChangeIndex) {
				continue
			}

			w.logger.Debug().
				Uint64("old", w.lastChangeIndex).
				Uint64("new", deployments[i].ModifyIndex).
				Str("deployment", deployments[i].ID).
				Msg("deployment modify index has changed is greater than last recorded")

			maxFound = watcher.MaxFound(deployments[i].ModifyIndex, maxFound)
			updateChan <- deployments[i]
		}

		// Update the Nomad API wait index to start long polling from the correct point and update
		// our recorded lastChangeIndex so we have the correct point to use during the next API
		// return.
		q.WaitIndex = meta.LastIndex
		w.lastChangeIndex = maxFound
	}
}
