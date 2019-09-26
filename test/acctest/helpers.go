package acctest

import (
	"fmt"
)

// CleanupPurgeJob is a cleanup func to purge the TestCase job from Nomad
func CleanupPurgeJob(s *TestState) error {
	_, _, err := s.Nomad.Jobs().Deregister(s.JobName, true, nil)
	// TODO: wait for action to complete
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
