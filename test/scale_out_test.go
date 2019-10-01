package test

import (
	"fmt"
	"os"
	"testing"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/api"
	"github.com/jrasell/sherpa/test/acctest"
)

const (
	testScaleOutGroupName1 = "sherpa-acctest-group-1"
	testScaleOutTaskName1  = "sherpa-acctest-task-1"
)

func TestScaleOut_singleTaskGroupCountSet(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") != "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{Enabled: true, MaxCount: 5, MinCount: 1}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, testScaleOutGroupName1, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testScaleOutGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}
					if policy1.MinCount != 1 {
						return fmt.Errorf("expected policy %s/%s to match the MinCount", s.JobName, testMetaGroupName1)
					}
					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildScaleOutTestJob(s.JobName), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					resp, err := s.Sherpa.Scale().JobGroupOut(s.JobName, testScaleOutGroupName1, 2)
					if err != nil {
						return err
					}

					if resp == nil {
						return fmt.Errorf("expected non-nil scale out response")
					}

					if _, err = s.Sherpa.Scale().Info(resp.ID.String()); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Runner: acctest.CheckTaskGroupCount(testScaleOutGroupName1, 3),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func TestScaleOut_singleTaskGroupCountSetTooHigh(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") != "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{Enabled: true, MaxCount: 2, MinCount: 1}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, testScaleOutGroupName1, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testScaleOutGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 2 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}
					if policy1.MinCount != 1 {
						return fmt.Errorf("expected policy %s/%s to match the MinCount", s.JobName, testMetaGroupName1)
					}
					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildScaleOutTestJob(s.JobName), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Scale().JobGroupOut(s.JobName, testScaleOutGroupName1, 10)
					if err != nil {
						return err
					}
					return nil
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 304:"),
			},
			{
				Runner: acctest.CheckTaskGroupCount(testScaleOutGroupName1, 1),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func buildScaleOutTestJob(name string) *nomad.Job {
	j := acctest.BuildBaseTestJob(name)
	j.TaskGroups = append(j.TaskGroups, acctest.BuildBaseTaskGroup(testScaleOutGroupName1, testScaleOutTaskName1))
	return j
}
