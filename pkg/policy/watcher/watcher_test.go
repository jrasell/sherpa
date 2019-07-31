package watcher

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMetaWatcher_indexHasChange(t *testing.T) {
	watcher := NewMetaWatcher(zerolog.Logger{}, nil, nil)

	testCases := []struct {
		newValue       uint64
		oldValue       uint64
		expectedReturn bool
	}{
		{
			newValue:       13,
			oldValue:       7,
			expectedReturn: true,
		},
		{
			newValue:       13696,
			oldValue:       13696,
			expectedReturn: false,
		},
		{
			newValue:       7,
			oldValue:       13,
			expectedReturn: false,
		},
	}

	for _, tc := range testCases {
		res := watcher.indexHasChange(tc.newValue, tc.oldValue)
		assert.Equal(t, tc.expectedReturn, res)
	}
}

func TestMetaWatcher_maxFound(t *testing.T) {
	watcher := NewMetaWatcher(zerolog.Logger{}, nil, nil)

	testCases := []struct {
		newValue       uint64
		oldValue       uint64
		expectedReturn uint64
	}{
		{
			newValue:       13,
			oldValue:       7,
			expectedReturn: 13,
		},
		{
			newValue:       13696,
			oldValue:       13696,
			expectedReturn: 13696,
		},
		{
			newValue:       7,
			oldValue:       13,
			expectedReturn: 13,
		},
		{
			newValue:       1,
			oldValue:       0,
			expectedReturn: 1,
		},
	}

	for _, tc := range testCases {
		res := watcher.maxFound(tc.newValue, tc.oldValue)
		assert.Equal(t, tc.expectedReturn, res)
	}
}
