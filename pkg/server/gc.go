package server

import "time"

func (h *HTTPServer) runGarbageCollectionLoop() {
	h.logger.Info().Msg("started scaling state garbage collector handler")

	t := time.NewTicker(time.Minute * 10)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			h.logger.Debug().Msg("triggering internal run of state garbage collection")
			h.stateBackend.RunGarbageCollection()
		}
	}
}
