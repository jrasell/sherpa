package scale

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ScaleConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterScaleConfig(fakeCMD)

	cfg := GetScaleConfig()
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, "", cfg.GroupName)
}
