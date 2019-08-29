package helper

import "time"

// UnixNanoToHumanUTC is used to convert the internally used UnixNano timestamp to a human readable
// output.
func UnixNanoToHumanUTC(t int64) time.Time {
	return time.Unix(0, t).UTC().Round(10 * time.Millisecond)
}
