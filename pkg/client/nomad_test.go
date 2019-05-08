package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewNomadClient(t *testing.T) {
	addr := "http://nomad.jrasell.system:4646"

	err := os.Setenv("NOMAD_ADDR", addr)
	assert.Nil(t, err)

	client, err := NewNomadClient()
	assert.Nil(t, err)

	nAddr := client.Address()
	assert.Equal(t, addr, nAddr)
}
