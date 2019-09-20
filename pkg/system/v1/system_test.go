package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/jrasell/sherpa/pkg/client"
	"github.com/jrasell/sherpa/pkg/config/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSystem_GetHealth(t *testing.T) {
	s := NewSystemServer(zerolog.Logger{}, nil, nil, nil, nil)

	r := httptest.NewRequest("GET", "http://jrasell.com/v1/system/health", nil)
	w := httptest.NewRecorder()
	s.GetHealth(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), defaultHealthResp)
}

func TestSystem_GetInfo(t *testing.T) {
	testCases := []struct {
		systemServerConfig *server.Config
		expectedRespCode   int
		expectedRespBody   string
	}{
		{
			systemServerConfig: &server.Config{APIPolicyEngine: true, ConsulStorageBackend: true},
			expectedRespCode:   200,
			expectedRespBody:   "{\"NomadAddress\":\"http://127.0.0.1:4646\",\"PolicyEngine\":\"Sherpa API\",\"StorageBackend\":\"Consul\",\"InternalAutoScalingEngine\":false,\"StrictPolicyChecking\":false}",
		},
		{
			systemServerConfig: &server.Config{NomadMetaPolicyEngine: true},
			expectedRespCode:   200,
			expectedRespBody:   "{\"NomadAddress\":\"http://127.0.0.1:4646\",\"PolicyEngine\":\"Nomad Job Group Meta\",\"StorageBackend\":\"In Memory\",\"InternalAutoScalingEngine\":false,\"StrictPolicyChecking\":false}",
		},
		{
			systemServerConfig: &server.Config{APIPolicyEngine: true, InternalAutoScaler: true},
			expectedRespCode:   200,
			expectedRespBody:   "{\"NomadAddress\":\"http://127.0.0.1:4646\",\"PolicyEngine\":\"Sherpa API\",\"StorageBackend\":\"In Memory\",\"InternalAutoScalingEngine\":true,\"StrictPolicyChecking\":false}",
		},
		{
			systemServerConfig: &server.Config{},
			expectedRespCode:   200,
			expectedRespBody:   "{\"NomadAddress\":\"http://127.0.0.1:4646\",\"PolicyEngine\":\"Disabled\",\"StorageBackend\":\"In Memory\",\"InternalAutoScalingEngine\":false,\"StrictPolicyChecking\":false}",
		},
	}

	nomadClient, _ := client.NewNomadClient()

	for _, tc := range testCases {
		r := httptest.NewRequest("GET", "http://jrasell.com/v1/system/info", nil)
		w := httptest.NewRecorder()

		s := NewSystemServer(zerolog.Logger{}, nomadClient, tc.systemServerConfig, nil, nil)
		s.GetInfo(w, r)

		assert.Equal(t, tc.expectedRespCode, w.Code)
		assert.Equal(t, tc.expectedRespBody, w.Body.String())
	}
}
