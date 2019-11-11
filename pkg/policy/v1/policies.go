package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Policy struct {
	logger  zerolog.Logger
	backend backend.PolicyBackend
}

func NewPolicyServer(l zerolog.Logger, backend backend.PolicyBackend) *Policy {
	return &Policy{logger: l, backend: backend}
}

func (p *Policy) GetJobPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := p.backend.GetPolicies()
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to call policy backend")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(policies)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to format HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, bytes, http.StatusOK)
}

func (p *Policy) GetJobPolicy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	job := vars["job_id"]

	policies, err := p.backend.GetJobPolicy(job)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to call policy backend")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if policies == nil {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(policies)
	if err != nil {
		p.logger.Error().Err(err).Msg(readBodyFailureMsg)
		http.Error(w, readBodyFailureMsg, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func (p *Policy) GetJobGroupPolicy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	job := vars["job_id"]
	group := vars["group"]

	gPolicy, err := p.backend.GetJobGroupPolicy(job, group)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to call policy backend")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if gPolicy == nil || gPolicy == (&policy.GroupScalingPolicy{}) {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(gPolicy)
	if err != nil {
		p.logger.Error().Err(err).Msg(marshalRespFailureMsg)
		http.Error(w, marshalRespFailureMsg, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func (p *Policy) PutJobPolicy(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.logger.Error().Msg(readBodyFailureMsg)
		http.Error(w, readBodyFailureMsg, http.StatusInternalServerError)
		return
	}

	jobPolicy, err := decodeJobPolicyReqBodyAndValidate(b)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to decode request body")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	job := vars["job_id"]

	if err := p.backend.PutJobPolicy(job, jobPolicy); err != nil {
		p.logger.Error().Err(err).Msg("failed to call policy backend")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *Policy) PutJobGroupPolicy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	job := vars["job_id"]
	group := vars["group"]

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.logger.Error().Msg(readBodyFailureMsg)
		http.Error(w, readBodyFailureMsg, http.StatusInternalServerError)
		return
	}

	groupPolicy, err := decodeGroupPolicyReqBodyAndValidate(b)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to decode request body")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := p.backend.PutJobGroupPolicy(job, group, groupPolicy); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *Policy) DeleteJobGroupPolicy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	job := vars["job_id"]
	group := vars["group"]

	if err := p.backend.DeleteJobGroupPolicy(job, group); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *Policy) DeleteJobPolicy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	job := vars["job_id"]

	if err := p.backend.DeleteJobPolicy(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSONResponse(w http.ResponseWriter, bytes []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
	}
}

func decodeGroupPolicyReqBodyAndValidate(body []byte) (*policy.GroupScalingPolicy, error) {
	p := &policy.GroupScalingPolicy{}

	if err := json.Unmarshal(body, p); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal request body")
	}

	if err := p.Validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate policy document")
	}
	return p.MergeWithDefaults(), nil
}

func decodeJobPolicyReqBodyAndValidate(body []byte) (map[string]*policy.GroupScalingPolicy, error) {
	p := make(map[string]*policy.GroupScalingPolicy)

	if err := json.Unmarshal(body, &p); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal request body")
	}

	for _, pol := range p {
		if err := pol.Validate(); err != nil {
			return nil, errors.Wrap(err, "failed to validate policy document")
		}
		pol = pol.MergeWithDefaults()
	}

	return p, nil
}
