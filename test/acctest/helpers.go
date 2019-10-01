package acctest

import (
	"fmt"
	"time"

	nomad "github.com/hashicorp/nomad/api"
)

func StringToPointer(s string) *string { return &s }
func IntToPointer(i int) *int          { return &i }

// CleanupPurgeJob is a cleanup func to purge the TestCase job from Nomad
func CleanupPurgeJob(s *TestState) error {
	_, _, err := s.Nomad.Jobs().Deregister(s.JobName, true, nil)
	return err
}

// CleanupSherpaPolicy is a CleanupFunc to remove a single job policy
func CleanupSherpaPolicy(s *TestState) error {
	return s.Sherpa.Policies().DeleteJobPolicy(s.JobName)
}

// CheckErrEqual is a CheckErr func to test if an error message is as expected
func CheckErrEqual(expected string) func(error) bool {
	return func(err error) bool {
		return expected == err.Error()
	}
}

// CheckJobReachesStatus performs a check, with a timeout that the test job reaches to desired
// status.
func CheckJobReachesStatus(s *TestState, status string) error {
	timeout := time.After(30 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout reached on checking that job %s reaches status %s", s.JobName, status)
		case <-tick:
			j, _, err := s.Nomad.Jobs().Info(s.JobName, nil)
			if err != nil {
				return err
			}
			if *j.Status == status {
				return nil
			}
		}
	}
}

// CheckDeploymentStatus is a TestStateFunc to check if the latest deployment of
// the TestCase job in Nomad matches the desired status
func CheckDeploymentStatus(status string) TestStateFunc {
	return func(s *TestState) error {
		deploy, _, err := s.Nomad.Jobs().LatestDeployment(s.JobName, nil)
		if err != nil {
			return err
		}

		if deploy.Status != status {
			return fmt.Errorf("deployment %s is in status '%s', expected '%s'", deploy.ID, deploy.Status, status)
		}

		return nil
	}
}

// CheckTaskGroupCount is a TestStateFunc to check a TaskGroup count matches the expected count
func CheckTaskGroupCount(groupName string, count int) TestStateFunc {
	return func(s *TestState) error {
		job, _, err := s.Nomad.Jobs().Info(s.JobName, nil)
		if err != nil {
			return err
		}

		for _, group := range job.TaskGroups {
			if groupName == *group.Name {
				if *group.Count == count {
					return nil
				}

				return fmt.Errorf("task group %s count is %d, expected %d", groupName, *group.Count, count)
			}
		}

		return fmt.Errorf("unable to find task group %s", groupName)
	}
}

func BuildBaseTestJob(name string) *nomad.Job {
	return &nomad.Job{
		ID:          StringToPointer(name),
		Name:        StringToPointer(name),
		Datacenters: []string{"dc1"},
	}
}

func BuildBaseTaskGroup(group, task string) *nomad.TaskGroup {
	return &nomad.TaskGroup{
		Name: StringToPointer(group),
		Tasks: []*nomad.Task{{
			Name:   task,
			Driver: "docker",
			Config: map[string]interface{}{"image": "redis:3.2"},
			Resources: &nomad.Resources{
				CPU:      IntToPointer(500),
				MemoryMB: IntToPointer(256),
			},
		}},
	}
}
