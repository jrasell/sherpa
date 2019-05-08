package v1

import (
	"net/http"

	"github.com/hashicorp/nomad/api"
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2/json"
)

const (
	headerKeyContentType       = "Content-Type"
	headerValueContentTypeJSON = "application/json; charset=utf-8"

	defaultHealthResp          = "{\"status\":\"ok\"}"
	defaultAPIPolicyResp       = "Sherpa API"
	defaultMetaPolicyResp      = "Nomad Job Group Meta"
	defaultDisabledPolicyResp  = "Disabled"
	defaultPolicyBackend       = "In Memory"
	defaultPolicyBackendConsul = "Consul"
)

type System struct {
	logger zerolog.Logger
	nomad  *api.Client
	server *serverCfg.Config
}

type SystemInfoResp struct {
	NomadAddress              string
	PolicyEngine              string
	PolicyStorageBackend      string
	InternalAutoScalingEngine bool
	StrictPolicyChecking      bool
}

func NewSystemServer(l zerolog.Logger, nomad *api.Client, server *serverCfg.Config) *System {
	return &System{
		logger: l,
		nomad:  nomad,
		server: server,
	}
}

func (h *System) GetHealth(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, []byte(defaultHealthResp), http.StatusOK)
}

func (h *System) GetInfo(w http.ResponseWriter, r *http.Request) {
	resp := &SystemInfoResp{
		NomadAddress:              h.nomad.Address(),
		StrictPolicyChecking:      h.server.StrictPolicyChecking,
		InternalAutoScalingEngine: h.server.InternalAutoScaler,
		PolicyEngine:              defaultDisabledPolicyResp,
		PolicyStorageBackend:      defaultPolicyBackend,
	}

	if h.server.ConsulStorageBackend {
		resp.PolicyStorageBackend = defaultPolicyBackendConsul
	}

	if h.server.APIPolicyEngine {
		resp.PolicyEngine = defaultAPIPolicyResp
	}

	if h.server.NomadMetaPolicyEngine {
		resp.PolicyEngine = defaultMetaPolicyResp
	}

	out, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal HTTP response")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, out, http.StatusOK)
}

func writeJSONResponse(w http.ResponseWriter, bytes []byte, statusCode int) {
	w.Header().Set(headerKeyContentType, headerValueContentTypeJSON)
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
	}
}
