package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ServerConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterConfig(fakeCMD)

	cfg := GetConfig()
	assert.Equal(t, configKeyBindAddrDefault, cfg.Bind)
	assert.Equal(t, uint16(configKeyBindPortDefault), cfg.Port)
	assert.Equal(t, true, cfg.APIPolicyEngine)
	assert.Equal(t, false, cfg.NomadMetaPolicyEngine)
	assert.Equal(t, true, cfg.StrictPolicyChecking)
	assert.Equal(t, false, cfg.InternalAutoScaler)
	assert.Equal(t, configKeyStorageBackendConsulPathDefault, cfg.ConsulStorageBackendPath)
	assert.Equal(t, configKeyAutoscalerThreadNumberDefault, cfg.InternalAutoScalerNumThreads)
}
