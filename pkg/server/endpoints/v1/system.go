package v1

import (
	"encoding/json"
	"net/http"

	metrics "github.com/armon/go-metrics"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/nomad/api"
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
	"github.com/jrasell/sherpa/pkg/server/cluster"
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

type SystemServer struct {
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

func NewSystemServer(l zerolog.Logger, nomad *api.Client, server *serverCfg.Config, tel *metrics.InmemSink, mem *cluster.Member) *SystemServer {
	return &SystemServer{
		logger:    l,
		member:    mem,
		nomad:     nomad,
		server:    server,
		telemetry: tel,
	}
}

func (s *SystemServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, []byte(defaultHealthResp))
}

func (s *SystemServer) GetInfo(w http.ResponseWriter, r *http.Request) {
	resp := &SystemInfoResp{
		NomadAddress:              s.nomad.Address(),
		StrictPolicyChecking:      s.server.StrictPolicyChecking,
		InternalAutoScalingEngine: s.server.InternalAutoScaler,
		PolicyEngine:              defaultDisabledPolicyResp,
		StorageBackend:            defaultStorageBackend,
	}

	if s.server.ConsulStorageBackend {
		resp.StorageBackend = defaultStorageBackendConsul
	}

	if s.server.APIPolicyEngine {
		resp.PolicyEngine = defaultAPIPolicyResp
	}

	if s.server.NomadMetaPolicyEngine {
		resp.PolicyEngine = defaultMetaPolicyResp
	}

	out, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal HTTP response")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, out)
}

func (s *SystemServer) GetLeader(w http.ResponseWriter, r *http.Request) {

	// Pull the leadership information from the local member.
	l, addr, advAddr, err := s.member.Leader()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get leadership information")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SystemLeaderResp{
		IsSelf:               l,
		HAEnabled:            s.member.IsHA(),
		LeaderAddress:        addr,
		LeaderClusterAddress: advAddr,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, out)
}

func (s *SystemServer) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metricData, err := s.telemetry.DisplayMetrics(w, r)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get latest telemetry data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(metricData)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal HTTP response")
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
