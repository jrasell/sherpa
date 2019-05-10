package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_TLSConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterTLSConfig(fakeCMD)

	cfg := GetTLSConfig()
	assert.Equal(t, "", cfg.CertKeyPath)
	assert.Equal(t, "", cfg.CertPath)
}
