package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jrasell/sherpa/pkg/helper"
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

// ScaleConfig is a convenience for setting up the scale server. These objects are centrally built
// and passed to the server.
type ScaleConfig struct {
	Logger zerolog.Logger
	Policy policyBackend.PolicyBackend
	Scale  scale.Scale
	State  stateBackend.Backend
}

type scaleRequestBody struct {
	Meta map[string]string
}

func NewScaleServer(strict bool, cfg *ScaleConfig) *Scale {
	return &Scale{
		logger:         cfg.Logger,
		scaler:         cfg.Scale,
		policyBackend:  cfg.Policy,
		stateBackend:   cfg.State,
		strictChecking: strict,
	}
}

func (s *Scale) InJobGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job_id"]
	groupID := vars["group"]

	body, err := parseScaleRequestBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newReq := &scale.GroupReq{
		Direction: scale.DirectionIn,
		GroupName: groupID,
		Time:      helper.GenerateEventTimestamp(),
		Meta:      body.Meta,
	}

	if s.scaler.JobGroupIsDeploying(jobID, groupID) {
		s.logger.Info().
			Str("job", jobID).
			Str("group", groupID).
			Msg("job group is currently in deployment and cannot be scaled")
		http.Error(w, errJobGroupInDeployment.Error(), http.StatusForbidden)
		return
	}

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

	if newReq.GroupScalingPolicy != nil {
		cd, err := s.scaler.JobGroupIsInCooldown(jobID, groupID, pol.Cooldown, newReq.Time)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("job", jobID).
				Str("group", groupID).
				Msg("failed to check if job group is currently in scaling cooldown")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if cd {
			s.logger.Info().
				Str("job", jobID).
				Str("group", groupID).
				Msg(jobGroupInCooldownMsg)
			http.Error(w, jobGroupInCooldownMsg, http.StatusConflict)
			return
		}
	}

	newReq.Count, err = payloadOrPolicyCount(getCountFromQueryParam(r), pol, scale.DirectionIn)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to determine scale count based on payload and policy")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scaleResp, respCode, err := s.scaler.Trigger(jobID, []*scale.GroupReq{newReq}, state.SourceAPI)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to scale in Nomad job group")
		http.Error(w, err.Error(), respCode)
		return
	}

	if respCode == http.StatusNotFound {
		http.NotFound(w, r)
		return
	}

	if respCode == http.StatusNotModified {
		http.Error(w, errors.New("unable to scale job").Error(), http.StatusNotModified)
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

	writeJSONResponse(w, bytes, http.StatusCreated)
}

func (s *Scale) OutJobGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job_id"]
	groupID := vars["group"]

	body, err := parseScaleRequestBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newReq := &scale.GroupReq{
		Direction: scale.DirectionOut,
		GroupName: groupID,
		Time:      helper.GenerateEventTimestamp(),
		Meta:      body.Meta,
	}

	if s.scaler.JobGroupIsDeploying(jobID, groupID) {
		s.logger.Info().
			Str("job", jobID).
			Str("group", groupID).
			Msg("job group is currently in deployment and cannot be scaled")
		http.Error(w, errJobGroupInDeployment.Error(), http.StatusForbidden)
		return
	}

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

	if newReq.GroupScalingPolicy != nil {
		cd, err := s.scaler.JobGroupIsInCooldown(jobID, groupID, pol.Cooldown, newReq.Time)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("job", jobID).
				Str("group", groupID).
				Msg("failed to check if job group is currently in scaling cooldown")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if cd {
			s.logger.Info().
				Str("job", jobID).
				Str("group", groupID).
				Msg(jobGroupInCooldownMsg)
			http.Error(w, jobGroupInCooldownMsg, http.StatusConflict)
			return
		}
	}

	newReq.Count, err = payloadOrPolicyCount(getCountFromQueryParam(r), pol, scale.DirectionIn)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to determine scale count based on payload and policy")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scaleResp, respCode, err := s.scaler.Trigger(jobID, []*scale.GroupReq{newReq}, state.SourceAPI)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("job", jobID).
			Str("group", groupID).
			Msg("failed to scale out Nomad job group")
		http.Error(w, err.Error(), respCode)
		return
	}

	if respCode == http.StatusNotFound {
		http.NotFound(w, r)
		return
	}

	if respCode == http.StatusNotModified {
		http.Error(w, "unable to scale job", http.StatusNotModified)
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

	writeJSONResponse(w, bytes, http.StatusCreated)
}

func parseScaleRequestBody(r *http.Request) (*scaleRequestBody, error) {
	if r.ContentLength < 1 {
		empty := &scaleRequestBody{
			Meta: map[string]string{},
		}
		return empty, nil
	}

	var body scaleRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
