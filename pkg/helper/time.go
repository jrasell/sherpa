package helper

import "time"

// GenerateEventTimestamp is used to provide a consistent way to generate timestamps across all
// internal events. Any changes here will impact all event messages.
func GenerateEventTimestamp() int64 {
	return time.Now().UTC().UnixNano()
}
