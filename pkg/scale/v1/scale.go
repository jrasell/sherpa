package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Scale struct {
	logger         zerolog.Logger
	policyBackend  backend.PolicyBackend
	strictChecking bool
	scaler         scale.Scale
}

func NewScaleServer(l zerolog.Logger, strict bool, backend backend.PolicyBackend, c *api.Client) *Scale {
	return &Scale{
		logger:         l,
		scaler:         scale.NewScaler(c, l, strict),
		policyBackend:  backend,
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

	scaleResp, respCode, err := s.scaler.Trigger(jobID, []*scale.GroupReq{newReq})
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
		Str("warnings", scaleResp.Warnings).
		Str("job", jobID).
		Str("group", groupID).
		Msg("successfully scaled in Nomad job group")

	writeJSONResponse(w, []byte(fmt.Sprintf("{\"EvaluationID\":\"%s\"}", scaleResp.EvalID)), http.StatusOK)
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

	scaleResp, respCode, err := s.scaler.Trigger(jobID, []*scale.GroupReq{newReq})
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
		Str("warnings", scaleResp.Warnings).
		Str("job", jobID).
		Str("group", groupID).
		Msg("successfully scaled out Nomad job group")

	writeJSONResponse(w, []byte(fmt.Sprintf("{\"EvaluationID\":\"%s\"}", scaleResp.EvalID)), http.StatusOK)
}
