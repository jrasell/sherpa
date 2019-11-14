package scale

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ScaleStatusConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterScaleConfig(fakeCMD)

	cfg := GetScaleStatusConfig()
	assert.Equal(t, false, cfg.Latest)
}
