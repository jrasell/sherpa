package sherpa

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_SherpaConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterConfig(fakeCMD)
	assert.Equal(t, configKeySherpaAddrDefault, GetSherpaConfig().Addr)
}
