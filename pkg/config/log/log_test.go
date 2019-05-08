package log

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ServerConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterConfig(fakeCMD)

	cfg := GetConfig()
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "auto", cfg.LogFormat)
	assert.Equal(t, false, cfg.EnableDev)
	assert.Equal(t, false, cfg.UseColor)
}
