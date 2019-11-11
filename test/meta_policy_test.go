package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/test/acctest"
)

const (
	testMetaGroupName1 = "sherpa-acctest-group-1"
	testMetaTaskName1  = "sherpa-acctest-task-1"
	testMetaGroupName2 = "sherpa-acctest-group-2"
	testMetaTaskName2  = "sherpa-acctest-task-2"
)

type meta string

const (
	metaPartial  meta = "partial"
	metaAll      meta = "all"
	metaNone     meta = "none"
	metaExternal meta = "external"
)

func TestMetaPolicy_singleTaskGroupFullMetaRemoveAllMeta(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(false, s.JobName, metaAll), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(false, s.JobName, metaNone), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_singleTaskGroupFullMetaStopJob(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(false, s.JobName, metaAll), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Deregister(s.JobName, false, nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "dead")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_multiTaskGroupFullMetaRemoveAllMeta(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaAll), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}

					policy2, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName2)
					if err != nil {
						return err
					}

					if policy2.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName2)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaNone), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_multiTaskGroupFullMetaStopJob(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaAll), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}

					policy2, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName2)
					if err != nil {
						return err
					}

					if policy2.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName2)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Deregister(s.JobName, false, nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "dead")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_multiTaskGroupPartialMetaRemoveAllMeta(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaPartial), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}
					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName2)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaNone), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_multiTaskGroupPartialMetaStopJob(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaPartial), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}
					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName2)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Deregister(s.JobName, false, nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "dead")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_multiTaskGroupPartialMetaAddAllMeta(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaPartial), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}
					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName2)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(true, s.JobName, metaAll), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if policy1.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName1)
					}

					policy2, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName2)
					if err != nil {
						return err
					}

					if policy2.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, testMetaGroupName2)
					}

					return nil
				},
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func TestMetaPolicy_singleTaskGroupAllMetaExternalCheck(t *testing.T) {
	if os.Getenv("SHERPA_ACC_META") == "" {
		t.SkipNow()
	}

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					_, _, err := s.Nomad.Jobs().Register(buildMetaTestJob(false, s.JobName, metaExternal), nil)
					if err != nil {
						return err
					}
					return acctest.CheckJobReachesStatus(s, "running")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy1, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, testMetaGroupName1)
					if err != nil {
						return err
					}

					if len(policy1.ExternalChecks) != 1 {
						return fmt.Errorf("expected policy %s/%s to have 1 external check", s.JobName, testMetaGroupName1)
					}
					return nil
				},
			},
		},
		CleanupFuncs: []acctest.TestStateFunc{acctest.CleanupPurgeJob},
	})
}

func buildMetaTestJob(multiGroup bool, name string, metaType meta) *api.Job {
	j := acctest.BuildBaseTestJob(name)
	j.TaskGroups = append(j.TaskGroups, acctest.BuildBaseTaskGroup(testMetaGroupName1, testMetaTaskName1))

	if multiGroup {
		j.TaskGroups = append(j.TaskGroups, acctest.BuildBaseTaskGroup(testMetaGroupName2, testMetaTaskName2))
	}

	switch metaType {
	case metaAll:
		for i := range j.TaskGroups {
			j.TaskGroups[i].Meta = buildMetaBasic()
		}
	case metaPartial:
		j.TaskGroups[0].Meta = buildMetaBasic()
	case metaExternal:
		j.TaskGroups[0].Meta = buildMetaWithExternal()
	case metaNone:
	}

	return j
}

func buildMetaBasic() map[string]string {
	return map[string]string{"sherpa_enabled": "true", "sherpa_max_count": "5"}
}

func buildMetaWithExternal() map[string]string {
	return map[string]string{"sherpa_enabled": "true", "sherpa_max_count": "5", "sherpa_external_checks": "{\"ExternalChecks\":{\"prometheus_test\":{\"Enabled\":true,\"Provider\":\"prometheus\",\"Query\":\"job:nomad_redis_cache_memory:percentage\",\"ComparisonOperator\":\"less-than\",\"ComparisonValue\":30,\"Action\":\"scale-in\"}}}"}
}
