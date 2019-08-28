package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/nomad/api"
	policyBackend "github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/jrasell/sherpa/pkg/state"
	stateBackend "github.com/jrasell/sherpa/pkg/state/scale"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Scale struct {
	logger         zerolog.Logger
	policyBackend  policyBackend.PolicyBackend
	stateBackend   stateBackend.Backend
	strictChecking bool
	scaler         scale.Scale
}

func NewScaleServer(l zerolog.Logger, strict bool, backend policyBackend.PolicyBackend, state stateBackend.Backend, c *api.Client) *Scale {
	return &Scale{
		logger:         l,
		scaler:         scale.NewScaler(c, l, state, strict),
		policyBackend:  backend,
		stateBackend:   state,
		strictChecking: strict,
	}
}

func (s *Scale) InJobGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job_id"]
	groupID := vars["group"]

	newReq := &scale.GroupReq{Direction: scale.DirectionIn, GroupName: groupID}

	pol, err := s.policyBackend.GetJobGroupPolicy(jobID, groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.strictChecking && pol == nil {
		s.logger.Info().
			Str("job", jobID).
			Msg("strict checking enabled and job group does not have scaling policy")
		http.Error(w, errInternalScaleInNoPolicy.Error(), http.StatusForbidden)
		return
	}
	newReq.GroupScalingPolicy = pol

	newReq.Count, err = payloadOrPolicyCount(getCountFromQueryParam(r), pol, scale.DirectionIn)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to determine scaleResp count based on payload and policy")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scaleResp, respCode, err := s.scaler.Trigger(jobID, []*scale.GroupReq{newReq}, state.SourceAPI)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to scaleResp in Nomad job group")
		http.Error(w, err.Error(), respCode)
		return
	}

	if respCode == http.StatusNotFound {
		http.NotFound(w, r)
		return
	}

	if respCode == http.StatusNotModified {
		http.Error(w, errors.New("unable to scaleResp job").Error(), http.StatusNotModified)
		return
	}

	s.logger.Info().
		Str("job", jobID).
		Str("group", groupID).
		Msg("successfully scaled in Nomad job group")

	bytes, err := json.Marshal(scaleResp)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal scaling response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func (s *Scale) OutJobGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job_id"]
	groupID := vars["group"]

	newReq := &scale.GroupReq{Direction: scale.DirectionOut, GroupName: groupID}

	pol, err := s.policyBackend.GetJobGroupPolicy(jobID, groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.strictChecking && pol == nil {
		s.logger.Info().
			Str("job", jobID).
			Str("group", groupID).
			Msg("strict checking enabled and job group does not have scaling policy")
		http.Error(w, errInternalScaleOutNoPolicy.Error(), http.StatusForbidden)
		return
	}
	newReq.GroupScalingPolicy = pol

	newReq.Count, err = payloadOrPolicyCount(getCountFromQueryParam(r), pol, scale.DirectionIn)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to determine scaleResp count based on payload and policy")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scaleResp, respCode, err := s.scaler.Trigger(jobID, []*scale.GroupReq{newReq}, state.SourceAPI)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to scaleResp out Nomad job group")
		http.Error(w, err.Error(), respCode)
		return
	}

	if respCode == http.StatusNotFound {
		http.NotFound(w, r)
		return
	}

	if respCode == http.StatusNotModified {
		http.Error(w, "unable to scaleResp job", http.StatusNotModified)
		return
	}

	s.logger.Info().
		Str("job", jobID).
		Str("group", groupID).
		Msg("successfully scaled out Nomad job group")

	bytes, err := json.Marshal(scaleResp)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal scaling response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}
