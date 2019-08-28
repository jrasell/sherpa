package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_UnixNanoToHumanUTC(t *testing.T) {
	expectedOutput, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2019-08-22 14:06:40.11 +0000 UTC")
	actualOutput := UnixNanoToHumanUTC(1566482800109501000)

	assert.Equal(t, expectedOutput, actualOutput)
}
