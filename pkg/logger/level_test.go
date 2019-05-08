package logger

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_LevelString(t *testing.T) {
	testCases := []struct {
		level            Level
		expectedResponse string
	}{
		{level: LevelDebug, expectedResponse: "debug"},
		{level: LevelInfo, expectedResponse: "info"},
		{level: LevelWarn, expectedResponse: "warn"},
		{level: LevelError, expectedResponse: "error"},
		{level: LevelFatal, expectedResponse: "fatal"},
	}

	for _, tc := range testCases {
		res := tc.level.String()
		assert.Equal(t, tc.expectedResponse, res)
	}
}

func TestLogger_setLogLevel(t *testing.T) {
	testCases := []struct {
		level            string
		expectedResponse interface{}
	}{
		{level: "debug", expectedResponse: LevelDebug},
		{level: "info", expectedResponse: LevelInfo},
		{level: "warn", expectedResponse: LevelWarn},
		{level: "error", expectedResponse: LevelError},
		{level: "fatal", expectedResponse: LevelFatal},
		{level: "nuke", expectedResponse: fmt.Errorf("unsupported error level: %q (supported levels: %s)", "nuke", strings.Join(logLevelsStr(), " "))},
	}

	for _, tc := range testCases {
		level, err := setLogLevel(tc.level)
		if err != nil {
			assert.Equal(t, tc.expectedResponse, err)
			break
		}
		assert.Equal(t, tc.expectedResponse, level)
	}
}
