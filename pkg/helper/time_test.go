package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateEventTimestamp(t *testing.T) {
	testTime := GenerateEventTimestamp()

	expectedDigits := 19

	actualDigits := 0
	for testTime != 0 {
		testTime /= 10
		actualDigits++
	}

	assert.Equal(t, expectedDigits, actualDigits, "pkg/helper")
}
