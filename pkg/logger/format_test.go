package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerFormat_String(t *testing.T) {
	testCases := []struct {
		format           Format
		expectedResponse string
	}{
		{format: FormatAuto, expectedResponse: "auto"},
		{format: FormatZerolog, expectedResponse: "zerolog"},
		{format: FormatHuman, expectedResponse: "human"},
	}

	for _, tc := range testCases {
		res := tc.format.String()
		assert.Equal(t, tc.expectedResponse, res)
	}
}

func TestLogger_getLogFormat(t *testing.T) {
	testCases := []struct {
		format           string
		expectedResponse interface{}
		expectedError    error
	}{
		{format: "auto", expectedResponse: FormatAuto, expectedError: nil},
		{format: "json", expectedResponse: FormatZerolog, expectedError: nil},
		{format: "zerolog", expectedResponse: FormatZerolog, expectedError: nil},
		{format: "human", expectedResponse: FormatHuman, expectedError: nil},
		{format: "robot", expectedResponse: FormatAuto, expectedError: fmt.Errorf("unsupported log format: \"robot\"")},
	}

	for _, tc := range testCases {
		res, err := getLogFormat(tc.format)
		assert.Equal(t, tc.expectedError, err)
		assert.Equal(t, tc.expectedResponse, res)
	}
}
