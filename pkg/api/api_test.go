package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultConfig(t *testing.T) {
	expectedAddr := "http://127.0.0.1:8000"
	config := DefaultConfig()
	assert.Equal(t, expectedAddr, config.Address)
}

func Test_NewClient(t *testing.T) {
	testCases := []struct {
		inputConfig              *Config
		expectedReturnError      error
		expectedClientConfigAddr string
	}{
		{
			inputConfig:              &Config{},
			expectedReturnError:      nil,
			expectedClientConfigAddr: "http://127.0.0.1:8000",
		},
		{
			inputConfig:              &Config{Address: "sherpa.jrasell.com"},
			expectedReturnError:      nil,
			expectedClientConfigAddr: "sherpa.jrasell.com",
		},
	}

	for _, tc := range testCases {
		client, err := NewClient(tc.inputConfig)
		assert.Equal(t, tc.expectedClientConfigAddr, client.config.Address)

		if tc.expectedReturnError != nil {
			assert.EqualError(t, err, tc.expectedReturnError.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}
