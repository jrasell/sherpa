package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jrasell/sherpa/pkg/state"
)

func (s *Scale) StatusList(w http.ResponseWriter, r *http.Request) {
	if l := r.URL.Query().Get("latest"); l == "true" {
		s.statusListLatest(w)
		return
	}

	list, err := s.stateBackend.GetScalingEvents()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get scaling events from state")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(list)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal scaling state response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func (s *Scale) statusListLatest(w http.ResponseWriter) {
	list, err := s.stateBackend.GetLatestScalingEvents()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get latest scaling events from state")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := map[uuid.UUID]map[string]*state.ScalingEvent{}
	for jg, event := range list {
		out[event.ID] = map[string]*state.ScalingEvent{jg: event}
	}

	bytes, err := json.Marshal(out)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal latest scaling state response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func (s *Scale) StatusInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	scaleID, err := uuid.FromString(id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to convert scale ID query parameter to UUID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := s.stateBackend.GetScalingEvent(scaleID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get scaling event from state")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info == nil {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(info)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal scaling state response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}
