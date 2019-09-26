// Package acctest provides a small testing framework for Sherpa
package acctest

import (
	"fmt"
	"os"
	"testing"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
)

// TestCase is a single test of Sherpa
type TestCase struct {
	// Steps are ran in order stopping on failure
	Steps []TestStep

	// CleanupFunc is called at the end of the TestCase if set
	CleanupFunc TestStateFunc
}

// TestStep is a single step within a TestCase
type TestStep struct {
	// Runner is used to execute the step
	Runner TestStateFunc

	// ExpectErr allows Runner to fail, use CheckErr to confirm error is correct
	ExpectErr bool

	// CheckErr is called if Runner fails and ExpectErr is true
	CheckErr func(error) bool
}

// TestState is the configuration for the TestCase
type TestState struct {
	// JobName is a concatenation of "sherpa" and the test function name
	JobName string

	Sherpa *api.Client
	Nomad  *nomad.Client
}

// TestStateFunc provides a TestStep with access to the test state
type TestStateFunc func(*TestState) error

// ComposeTestStateFunc combines multiple TestStateFuncs into one
func ComposeTestStateFunc(f ...TestStateFunc) TestStateFunc {
	return func(s *TestState) error {
		for _, fun := range f {
			if err := fun(s); err != nil {
				return err
			}
		}

		return nil
	}
}

// Test executes a single TestCase
//
// Tests will be skipped if SHERPA_ACC is empty
func Test(t *testing.T, c TestCase) {
	if os.Getenv("SHERPA_ACC") == "" {
		t.SkipNow()
	}

	if len(c.Steps) < 1 {
		t.Fatal("must have at least one test step")
	}

	sherpa, err := newSherpaClient()
	if err != nil {
		t.Fatalf("failed to create Sherpa client: %s", err)
	}

	nomad, err := newNomadClient()
	if err != nil {
		t.Fatalf("failed to create Nomad client: %s", err)
	}

	state := &TestState{
		JobName: fmt.Sprintf("sherpa-%s", t.Name()),
		Sherpa:  sherpa,
		Nomad:   nomad,
	}

	for i, step := range c.Steps {
		stepNum := i + 1

		if step.Runner == nil {
			t.Errorf("step %d/%d does not have a Runner", stepNum, len(c.Steps))
			break
		}

		err = step.Runner(state)
		if err != nil {
			if !step.ExpectErr {
				t.Errorf("step %d/%d failed: %s", stepNum, len(c.Steps), err)
				break
			}

			if step.CheckErr != nil {
				ok := step.CheckErr(err)
				if !ok {
					t.Errorf("step %d/%d CheckErr failed: %s", stepNum, len(c.Steps), err)
					break
				}
			}
		}
	}

	if c.CleanupFunc != nil {
		err = c.CleanupFunc(state)
		if err != nil {
			t.Errorf("cleanup failed: %s", err)
		}
	}
}

// newNomadClient creates a Nomad API client configrable by NOMAD_
// env variables or returns an error if Nomad is in an unhealthy state
func newNomadClient() (*nomad.Client, error) {
	c, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		return nil, err
	}

	resp, err := c.Agent().Health()
	if err != nil {
		return nil, err
	}

	if !resp.Server.Ok || !resp.Client.Ok {
		return nil, fmt.Errorf("agent unhealthy")
	}

	return c, nil
}

// newNomadClient creates a Sherpa API client
func newSherpaClient() (*api.Client, error) {
	cfg := clientCfg.GetConfig()
	c, err := api.NewClient(api.DefaultConfig(&cfg))
	if err != nil {
		return nil, err
	}

	resp, err := c.System().Health()
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("Sherpa agent unhealthy, status: %s", resp.Status)
	}

	return c, nil
}
