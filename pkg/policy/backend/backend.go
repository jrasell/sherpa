package backend

import "github.com/jrasell/sherpa/pkg/policy"

// PolicyBackend is the interface required for a policy storage backend. A policy storage backend
// is used to durably store job scaling policies outside of Sherpa.
type PolicyBackend interface {
	// PutJobPolicy is used to insert or update the scaling policy for a job and all its associated
	// task groups.
	PutJobPolicy(string, map[string]*policy.GroupScalingPolicy) error

	// PutJobGroupPolicy is used to insert or update the scaling policy of a task group.
	PutJobGroupPolicy(string, string, *policy.GroupScalingPolicy) error

	// GetPolicies is used to retrieve all currently configured job scaling policies.
	GetPolicies() (map[string]map[string]*policy.GroupScalingPolicy, error)

	// GetJobPolicy retrieves the scaling policy for a job.
	GetJobPolicy(string) (map[string]*policy.GroupScalingPolicy, error)

	// GetJobGroupPolicy retrieves the scaling policy for a job task group.
	GetJobGroupPolicy(string, string) (*policy.GroupScalingPolicy, error)

	// DeleteJobPolicy is used to delete all task group scaling policies associated to the named
	// job.
	DeleteJobPolicy(string) error

	// DeleteJobGroupPolicy deletes the stored policy for a particular job group.
	DeleteJobGroupPolicy(string, string) error
}
