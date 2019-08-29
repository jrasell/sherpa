package memory

import (
	"sync"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
)

var _ backend.PolicyBackend = (*PolicyBackend)(nil)

type PolicyBackend struct {
	policies map[string]map[string]*policy.GroupScalingPolicy
	sync.RWMutex
}

func NewJobScalingPolicies() backend.PolicyBackend {
	return &PolicyBackend{
		policies: make(map[string]map[string]*policy.GroupScalingPolicy),
	}
}

func (p *PolicyBackend) GetPolicies() (map[string]map[string]*policy.GroupScalingPolicy, error) {
	p.RLock()
	val := p.policies
	p.RUnlock()
	return val, nil
}

func (p *PolicyBackend) GetJobPolicy(job string) (map[string]*policy.GroupScalingPolicy, error) {
	p.RLock()
	defer p.RUnlock()

	if val, ok := p.policies[job]; ok {
		return val, nil
	}
	return nil, nil
}

func (p *PolicyBackend) GetJobGroupPolicy(job, group string) (*policy.GroupScalingPolicy, error) {
	p.RLock()
	defer p.RUnlock()

	if val, ok := p.policies[job][group]; ok {
		return val, nil
	}
	return nil, nil
}

func (p *PolicyBackend) PutJobPolicy(job string, policies map[string]*policy.GroupScalingPolicy) error {
	p.Lock()
	defer p.Unlock()

	// A call to AddJobPolicy will overwrite the existing job policy, therefore here we initialise
	// the map entry.
	p.policies[job] = make(map[string]*policy.GroupScalingPolicy)

	for group, pol := range policies {
		p.policies[job][group] = pol
	}
	return nil
}

func (p *PolicyBackend) PutJobGroupPolicy(job, group string, policies *policy.GroupScalingPolicy) error {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.policies[job]; !ok {
		p.policies[job] = make(map[string]*policy.GroupScalingPolicy)
		p.policies[job][group] = policies
		return nil
	}

	p.policies[job][group] = policies
	return nil
}

func (p *PolicyBackend) DeleteJobGroupPolicy(job, group string) error {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.policies[job][group]; ok {
		delete(p.policies[job], group)
	}
	return nil
}

func (p *PolicyBackend) DeleteJobPolicy(job string) error {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.policies[job]; ok {
		delete(p.policies, job)
	}
	return nil
}
