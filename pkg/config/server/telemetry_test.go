package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_TelemetryConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterTelemetryConfig(fakeCMD)

	cfg := GetTelemetryConfig()
	assert.Equal(t, "", cfg.StatsiteAddr)
	assert.Equal(t, "", cfg.StatsdAddr)
}
