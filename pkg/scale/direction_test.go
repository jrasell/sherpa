package scale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirection_String(t *testing.T) {
	testCases := []struct {
		direction    Direction
		expectedResp string
	}{
		{direction: DirectionOut, expectedResp: "out"},
		{direction: DirectionIn, expectedResp: "in"},
	}

	for _, tc := range testCases {
		actualResp := tc.direction.String()
		assert.Equal(t, tc.expectedResp, actualResp)
	}
}
