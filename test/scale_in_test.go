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
	testScaleInGroupName1 = "sherpa-acctest-group-1"
	testScaleInTaskName1  = "sherpa-acctest-task-1"
)

func TestScaleIn_singleTaskGroupCountSet(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") != "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{Enabled: true, MaxCount: 5, MinCount: 1}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, testScaleInGroupName1, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testScaleInGroupName1)
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
					_, _, err := s.Nomad.Jobs().Register(buildScaleInTestJob(s.JobName), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					resp, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 2)
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
				Runner: acctest.CheckTaskGroupCount(testScaleInGroupName1, 1),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func TestScaleIn_singleTaskGroupCountSetTooLow(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") != "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{Enabled: true, MaxCount: 2, MinCount: 1}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, testScaleInGroupName1, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testScaleInGroupName1)
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
					_, _, err := s.Nomad.Jobs().Register(buildScaleInTestJob(s.JobName), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 10)
					if err != nil {
						return err
					}
					return nil
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 304:"),
			},
			{
				Runner: acctest.CheckTaskGroupCount(testScaleInGroupName1, 3),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func buildScaleInTestJob(name string) *nomad.Job {
	j := acctest.BuildBaseTestJob(name)
	j.TaskGroups = append(j.TaskGroups, acctest.BuildBaseTaskGroup(testScaleInGroupName1, testScaleInTaskName1))
	j.TaskGroups[0].Count = acctest.IntToPointer(3)
	return j
}
