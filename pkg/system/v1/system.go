package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/server/cluster"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/nomad/api"
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	headerKeyContentType       = "Content-Type"
	headerValueContentTypeJSON = "application/json; charset=utf-8"

	defaultHealthResp           = "{\"status\":\"ok\"}"
	defaultAPIPolicyResp        = "Sherpa API"
	defaultMetaPolicyResp       = "Nomad Job Group Meta"
	defaultDisabledPolicyResp   = "Disabled"
	defaultStorageBackend       = "In Memory"
	defaultStorageBackendConsul = "Consul"
)

type System struct {
	logger    zerolog.Logger
	member    *cluster.Member
	nomad     *api.Client
	server    *serverCfg.Config
	telemetry *metrics.InmemSink
}

type SystemInfoResp struct {
	NomadAddress              string
	PolicyEngine              string
	StorageBackend            string
	InternalAutoScalingEngine bool
	StrictPolicyChecking      bool
}

type SystemStatusResp struct {
	ClusterName string
	ClusterID   uuid.UUID
	HAEnabled   string
	Version     string
}

type SystemLeaderResp struct {
	IsSelf               bool
	HAEnabled            bool
	LeaderAddress        string
	LeaderClusterAddress string
}

func NewSystemServer(l zerolog.Logger, nomad *api.Client, server *serverCfg.Config, tel *metrics.InmemSink, mem *cluster.Member) *System {
	return &System{
		logger:    l,
		member:    mem,
		nomad:     nomad,
		server:    server,
		telemetry: tel,
	}
}

func (h *System) GetHealth(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, []byte(defaultHealthResp))
}

func (h *System) GetInfo(w http.ResponseWriter, r *http.Request) {
	resp := &SystemInfoResp{
		NomadAddress:              h.nomad.Address(),
		StrictPolicyChecking:      h.server.StrictPolicyChecking,
		InternalAutoScalingEngine: h.server.InternalAutoScaler,
		PolicyEngine:              defaultDisabledPolicyResp,
		StorageBackend:            defaultStorageBackend,
	}

	if h.server.ConsulStorageBackend {
		resp.StorageBackend = defaultStorageBackendConsul
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

	writeJSONResponse(w, out)
}

func (h *System) GetLeader(w http.ResponseWriter, r *http.Request) {

	// Pull the leadership information from the local member.
	l, addr, advAddr, err := h.member.Leader()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get leadership information")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SystemLeaderResp{
		IsSelf:               l,
		HAEnabled:            h.member.IsHA(),
		LeaderAddress:        addr,
		LeaderClusterAddress: advAddr,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, out)
}

func (h *System) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metricData, err := h.telemetry.DisplayMetrics(w, r)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get latest telemetry data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(metricData)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, out)
}

func writeJSONResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set(headerKeyContentType, headerValueContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
	}
}
