package autoscale

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_NomadConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterConfig(fakeCMD)

	cfg := GetConfig()
	assert.Equal(t, configKeyAutoscalingEvaluationIntervalDefault, cfg.AutoscalingEvalInt)
}
