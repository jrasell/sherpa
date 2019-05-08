package build

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

var (
	// GitDate is the date and time of the git commit used for a build (not the
	// build time)
	GitDate string

	// GitCommit is the SHA of the commit from which the platform was built.
	GitCommit string

	// GitBranch is the branch name from which the platform was built.
	GitBranch string

	// GitState indicates whether uncommitted changes were present at build
	// time.
	GitState string

	// ProjectName is the name of the complete system, used for config
	// directories etc
	ProjectName = "sherpa"

	// Version of the build
	Version string
)

func GetVersion() string {
	gitCommit := GitCommit
	if len(gitCommit) >= 8 {
		if _, err := hex.DecodeString(gitCommit); err == nil {
			gitCommit = gitCommit[:7]
		}
	}

	switch {
	case GitDate == "":
		GitDate = time.Now().UTC().String()
	default:
		if epoch, err := strconv.ParseUint(GitDate, 0, 64); err == nil {
			GitDate = time.Unix(int64(epoch), 0).UTC().String()
		}
	}

	version := ""
	if Version != "" {
		version = fmt.Sprintf("v%s", Version)
	}
	return fmt.Sprintf("%s\n\tDate: %s\n\tCommit: %s\n\tBranch: %s\n\tState: %s",
		version, GitDate, gitCommit, GitBranch, GitState)
}
