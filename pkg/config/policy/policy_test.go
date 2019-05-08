package policy

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_PolicyConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterConfig(fakeCMD)

	cfg := GetConfig()
	assert.Equal(t, "", cfg.GroupName)
}
