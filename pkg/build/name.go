package build

import (
	"os"
	"path"
)

var programName = ""

// ProgramName returns the name of the currently executing program. Perhaps misleadingly,
// this is set at runtime from cmd.Use of the root Cobra command.
func ProgramName() string {
	if programName == "" {
		programName = path.Base(os.Args[0])
	}

	return programName
}

// SetProgramName sets the name of the currently executing program. It should be called at
// the earliest opportunity (e.g. in PersistentPreRun of the root Cobra command).
func SetProgramName(name string) {
	programName = path.Base(name)
}
