package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ClusterConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterClusterConfig(fakeCMD)

	cfg := GetClusterConfig()
	assert.Equal(t, configKeyClusterAdvertiseAddrDefault, cfg.Addr)
	assert.Equal(t, "", cfg.Name)
}
