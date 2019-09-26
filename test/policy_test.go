package test

import (
	"fmt"
	"testing"

	"github.com/jrasell/sherpa/pkg/api"
	"github.com/jrasell/sherpa/test/acctest"
)

func TestPolicy_list(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					policies, err := s.Sherpa.Policies().List()
					if err != nil {
						return err
					}
					p := *policies

					if _, ok := p[s.JobName]; ok {
						return fmt.Errorf("expected policy %s to not exist", s.JobName)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{
						Enabled:  true,
						MaxCount: 5,
					}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, "group", policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policies, err := s.Sherpa.Policies().List()
					if err != nil {
						return err
					}
					p := *policies

					if _, ok := p[s.JobName]; !ok {
						return fmt.Errorf("expected policy %s to exist", s.JobName)
					}

					return nil
				},
			},
		},
		CleanupFunc: acctest.CleanupSherpaPolicy,
	})
}

func TestPolicy_readJob(t *testing.T) {
	groupName := "group"

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
					policy := &api.JobGroupPolicy{
						Enabled:  true,
						MaxCount: 5,
					}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, groupName, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policies, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					if err != nil {
						return err
					}
					p := *policies

					if _, ok := p[groupName]; !ok {
						return fmt.Errorf("expected policy %s/%s to exist", s.JobName, groupName)
					}

					return nil
				},
			},
		},
		CleanupFunc: acctest.CleanupSherpaPolicy,
	})
}

func TestPolicy_readJobGroup(t *testing.T) {
	groupName := "group"

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, groupName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{
						Enabled:  true,
						MaxCount: 5,
					}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, groupName, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, groupName)
					if err != nil {
						return err
					}

					if policy.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, groupName)
					}

					return nil
				},
			},
		},
		CleanupFunc: acctest.CleanupSherpaPolicy,
	})
}

func TestPolicy_write(t *testing.T) {
	groupName := "group"

	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: func(s *acctest.TestState) error {
					_, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, groupName)
					return err
				},
				ExpectErr: true,
				CheckErr:  acctest.CheckErrEqual("unexpected response code 404: 404 page not found"),
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{
						Enabled:  true,
						MaxCount: 5,
					}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, groupName, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, groupName)
					if err != nil {
						return err
					}

					if policy.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, groupName)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy := &api.JobGroupPolicy{
						Enabled:  true,
						MaxCount: 6,
					}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, groupName, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, groupName)
					if err != nil {
						return err
					}

					if policy.MaxCount != 6 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, groupName)
					}

					return nil
				},
			},
		},
		CleanupFunc: acctest.CleanupSherpaPolicy,
	})
}

func TestPolicy_deleteJobPolicy(t *testing.T) {
	groupName := "group"

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
					policy := &api.JobGroupPolicy{
						Enabled:  true,
						MaxCount: 5,
					}
					return s.Sherpa.Policies().WriteJobGroupPolicy(s.JobName, groupName, policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobGroupPolicy(s.JobName, groupName)
					if err != nil {
						return err
					}

					if policy.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to match the MaxCount", s.JobName, groupName)
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					return s.Sherpa.Policies().DeleteJobPolicy(s.JobName)
				},
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
		CleanupFunc: acctest.CleanupSherpaPolicy,
	})
}

func TestPolicy_deleteJobGroupPolicy(t *testing.T) {
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
					policy := map[string]*api.JobGroupPolicy{
						"group1": &api.JobGroupPolicy{
							Enabled:  true,
							MaxCount: 5,
						},
						"group2": &api.JobGroupPolicy{
							Enabled:  true,
							MaxCount: 10,
						},
					}
					return s.Sherpa.Policies().WriteJobPolicy(s.JobName, &policy)
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					if err != nil {
						return err
					}
					policyMap := *policy

					if p, ok := policyMap["group1"]; !ok || p.MaxCount != 5 {
						return fmt.Errorf("expected policy %s/%s to exist and match the MaxCount", s.JobName, "group1")
					}

					if p, ok := policyMap["group2"]; !ok || p.MaxCount != 10 {
						return fmt.Errorf("expected policy %s/%s to exist and match the MaxCount", s.JobName, "group2")
					}

					return nil
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					return s.Sherpa.Policies().DeleteJobGroupPolicy(s.JobName, "group1")
				},
			},
			{
				Runner: func(s *acctest.TestState) error {
					policy, err := s.Sherpa.Policies().ReadJobPolicy(s.JobName)
					if err != nil {
						return err
					}
					policyMap := *policy

					if _, ok := policyMap["group1"]; ok {
						return fmt.Errorf("expected policy %s/%s to not exist", s.JobName, "group1")
					}

					if p, ok := policyMap["group2"]; !ok || p.MaxCount != 10 {
						return fmt.Errorf("expected policy %s/%s to exist and match the MaxCount", s.JobName, "group2")
					}

					return nil
				},
			},
		},
		CleanupFunc: acctest.CleanupSherpaPolicy,
	})
}
