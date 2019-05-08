package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProgramName(t *testing.T) {
	programName = "SherpaProgramNameTest"
	returnProgramName := ProgramName()
	assert.Equal(t, programName, returnProgramName)
}

func Test_SetProgramName(t *testing.T) {
	SetProgramName("SherpaSetProgramNameTest")
	assert.Equal(t, programName, "SherpaSetProgramNameTest")
}
