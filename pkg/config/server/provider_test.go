package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_MetricProviderConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterMetricProviderConfig(fakeCMD)

	cfg := GetMetricProviderConfig()
	assert.Nil(t, cfg.Prometheus)
}
