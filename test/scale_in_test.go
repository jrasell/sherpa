package test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
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
					resp, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 2, nil)
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
					_, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 10, nil)
					if err != nil {
						return err
					}
					return nil
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 409: scaling action will break job group minimum threshold"),
			},
			{
				Runner: acctest.CheckTaskGroupCount(testScaleInGroupName1, 3),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func TestScaleIn_singleTaskGroupPolicyDisabled(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") != "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{Enabled: false, MaxCount: 2, MinCount: 1}
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
					_, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 10, nil)
					if err != nil {
						return err
					}
					return nil
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 409: job group scaling policy is currently disabled"),
			},
			{
				Runner: acctest.CheckTaskGroupCount(testScaleInGroupName1, 3),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func TestScaleIn_singleTaskGroupMeta(t *testing.T) {
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
					_, _, err := s.Nomad.Jobs().Register(buildScaleInTestJob(s.JobName), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					meta := map[string]string{
						"test-name": s.JobName,
					}
					resp, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 1, meta)
					if err != nil {
						return err
					}

					event, err := s.Sherpa.Scale().Info(resp.ID.String())
					if err != nil {
						return err
					}

					e, ok := event[s.JobName+":"+testScaleInGroupName1]
					if !ok {
						return errors.New("scaling event group not found")
					}

					if !reflect.DeepEqual(meta, e.Meta) {
						return errors.New("meta in event does not equal input")
					}

					return nil
				},
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupSherpaPolicy, acctest.CleanupPurgeJob},
	})
}

func TestScaleIn_singleTaskGroupWithHighCooldown(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") != "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{Enabled: true, Cooldown: 600, MinCount: 1}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, testScaleInGroupName1, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testScaleInGroupName1)
					if err != nil {
						return err
					}

					if policy1.Cooldown != 600 {
						return fmt.Errorf("expected policy %s/%s to match the Cooldown", s.JobName, testMetaGroupName1)
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
					resp, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 2, nil)
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
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Scale().JobGroupIn(s.JobName, testScaleInGroupName1, 1, nil)
					if err != nil {
						return err
					}
					return nil
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 409: job group is currently in scaling cooldown"),
			},
			{
				Runner: acctest.CheckTaskGroupCount(testScaleInGroupName1, 1),
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
