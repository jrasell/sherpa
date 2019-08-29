package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FormatList(t *testing.T) {
	input := []string{
		"ID|Job:Group|Status|Time",
		"4db95964-8a45-415e-b9ac-3ec3a9748e00|example1:cache|Completed|2019-08-23 09:55:51.009 +0000 UTC"}

	actualOutput := FormatList(input)
	expectedOutput := `ID                                    Job:Group       Status     Time
4db95964-8a45-415e-b9ac-3ec3a9748e00  example1:cache  Completed  2019-08-23 09:55:51.009 +0000 UTC`

	assert.Equal(t, expectedOutput, actualOutput)
}

func Test_FormatKV(t *testing.T) {
	input := []string{
		"ID|a11a6b4c-795e-4cd5-9fb1-56f7b9725875",
		"EvalID|a0a42e8d-4b96-9613-9a49-63d46a480bd9",
		"Status|Completed",
		"Source|InternalAutoscaler",
		"Time|2019-08-23 09:55:51.009 +0000 UTC",
	}

	actualOutput := FormatKV(input)
	expectedOutput := `ID     = a11a6b4c-795e-4cd5-9fb1-56f7b9725875
EvalID = a0a42e8d-4b96-9613-9a49-63d46a480bd9
Status = Completed
Source = InternalAutoscaler
Time   = 2019-08-23 09:55:51.009 +0000 UTC`

	assert.Equal(t, expectedOutput, actualOutput)
}
