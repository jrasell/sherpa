package api

import (
	"fmt"
)

type Policies struct {
	client *Client
}

func (c *Client) Policies() *Policies {
	return &Policies{client: c}
}

type JobGroupPolicy struct {
	Enabled       bool
	MaxCount      int
	MinCount      int
	ScaleOutCount int
	ScaleInCount  int

	ScaleOutCPUPercentageThreshold    int
	ScaleOutMemoryPercentageThreshold int
	ScaleInCPUPercentageThreshold     int
	ScaleInMemoryPercentageThreshold  int
}

func (p *Policies) List() (*map[string]map[string]*JobGroupPolicy, error) {
	var resp map[string]map[string]*JobGroupPolicy
	err := p.client.get("/v1/policies", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *Policies) ReadJobPolicy(job string) (*map[string]*JobGroupPolicy, error) {
	var resp map[string]*JobGroupPolicy
	err := p.client.get("/v1/policy/"+job, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *Policies) ReadJobGroupPolicy(job, group string) (*JobGroupPolicy, error) {
	var resp JobGroupPolicy

	path := fmt.Sprintf("/v1/policy/%s/%s", job, group)

	err := p.client.get(path, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *Policies) WriteJobPolicy(job string, policy *map[string]*JobGroupPolicy) error {
	return p.client.post("/v1/policy/"+job, policy, nil)
}

func (p *Policies) WriteJobGroupPolicy(job, group string, policy *JobGroupPolicy) error {
	path := fmt.Sprintf("/v1/policy/%s/%s", job, group)

	return p.client.post(path, policy, nil)
}

func (p *Policies) DeleteJobPolicy(job string) error {
	return p.client.delete("/v1/policy/"+job, nil)
}

func (p *Policies) DeleteJobGroupPolicy(job, group string) error {
	path := fmt.Sprintf("/v1/policy/%s/%s", job, group)
	return p.client.delete(path, nil)
}
