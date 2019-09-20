package server

import "time"

func (h *HTTPServer) runGarbageCollectionLoop() {
	h.logger.Info().Msg("started scaling state garbage collector handler")

	h.gcIsRunning = true

	t := time.NewTicker(time.Minute * 10)
	defer t.Stop()

	for {
		select {
		case <-h.stopChan:
			h.logger.Info().Msg("shutting down state garbage collection handler")
			h.gcIsRunning = false
			return
		case <-t.C:
			h.logger.Debug().Msg("triggering internal run of state garbage collection")
			h.stateBackend.RunGarbageCollection()
		}
	}
}
