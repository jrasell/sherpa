package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_DebugConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterDebugConfig(fakeCMD)
	cfg := GetClusterConfig()
	assert.False(t, GetDebugEnabled(), cfg.Addr)
}
