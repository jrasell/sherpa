package api

import (
	"testing"

	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	"github.com/stretchr/testify/assert"
)

func Test_DefaultConfig(t *testing.T) {
	testCases := []struct {
		inputConfig             *clientCfg.Config
		expectedAddrReturn      string
		expectedTLSConfigReturn *TLSConfig
	}{
		{
			inputConfig:             &clientCfg.Config{},
			expectedAddrReturn:      "http://127.0.0.1:8000",
			expectedTLSConfigReturn: &TLSConfig{},
		},
	}

	for _, tc := range testCases {
		actualReturn := DefaultConfig(tc.inputConfig)

		assert.Equal(t, tc.expectedAddrReturn, actualReturn.Address)
		assert.Equal(t, tc.expectedTLSConfigReturn, actualReturn.TLSConfig)
	}
}
