package memory

import (
	"testing"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/stretchr/testify/assert"
)

func TestPolicyBackend_Memory(t *testing.T) {
	newBackend := NewJobScalingPolicies()

	// Test putting a job group policy.
	sherpaGroup1 := generateTestPolicy()
	err := newBackend.PutJobGroupPolicy("sherpa-test-job-1", "sherpa-test-group-1", sherpaGroup1)
	assert.Nil(t, err)

	// Test reading the job group back.
	readSherpaGroup1, err := newBackend.GetJobGroupPolicy("sherpa-test-job-1", "sherpa-test-group-1")
	assert.Nil(t, err)
	assert.Equal(t, sherpaGroup1, readSherpaGroup1)

	// Test adding a second job group.
	sherpaGroup2 := generateTestPolicy()
	err = newBackend.PutJobGroupPolicy("sherpa-test-job-1", "sherpa-test-group-2", sherpaGroup2)
	assert.Nil(t, err)

	// Test reading the whole job back.
	expectedJob1 := map[string]*policy.GroupScalingPolicy{
		"sherpa-test-group-1": generateTestPolicy(),
		"sherpa-test-group-2": generateTestPolicy(),
	}
	readSherpaJob1, err := newBackend.GetJobPolicy("sherpa-test-job-1")
	assert.Nil(t, err)
	assert.Equal(t, expectedJob1, readSherpaJob1)

	// Test deleting a job group.
	err = newBackend.DeleteJobGroupPolicy("sherpa-test-job-1", "sherpa-test-group-1")
	assert.Nil(t, err)

	// Read the job back and ensure the record has been deleted.
	expectedJob2 := map[string]*policy.GroupScalingPolicy{"sherpa-test-group-2": generateTestPolicy()}
	readSherpaJob2, err := newBackend.GetJobPolicy("sherpa-test-job-1")
	assert.Nil(t, err)
	assert.Equal(t, expectedJob2, readSherpaJob2)

	// Test deleting a job policy.
	err = newBackend.DeleteJobPolicy("sherpa-test-job-1")
	assert.Nil(t, err)

	// Read the policies and check all are gone.
	expectedSherpaPolicies1 := map[string]map[string]*policy.GroupScalingPolicy{}
	sherpaPolicies1, err := newBackend.GetPolicies()
	assert.Nil(t, err)
	assert.Equal(t, expectedSherpaPolicies1, sherpaPolicies1)

	// Test putting a whole job policy.
	putSherpaJob1 := map[string]*policy.GroupScalingPolicy{
		"sherpa-test-group-3": generateTestPolicy(),
		"sherpa-test-group-4": generateTestPolicy(),
	}
	err = newBackend.PutJobPolicy("sherpa-test-2", putSherpaJob1)
	assert.Nil(t, err)

	// Read the job back and ensure it has been actually added.
	readSherpaJob3, err := newBackend.GetJobPolicy("sherpa-test-2")
	assert.Nil(t, err)
	assert.Equal(t, putSherpaJob1, readSherpaJob3)
}

func generateTestPolicy() *policy.GroupScalingPolicy {
	return &policy.GroupScalingPolicy{
		Enabled:                           true,
		MinCount:                          1,
		MaxCount:                          10,
		ScaleInCount:                      1,
		ScaleOutCount:                     2,
		ScaleOutCPUPercentageThreshold:    80,
		ScaleInCPUPercentageThreshold:     20,
		ScaleOutMemoryPercentageThreshold: 80,
		ScaleInMemoryPercentageThreshold:  20,
	}
}
