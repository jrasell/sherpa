package api

import metrics "github.com/armon/go-metrics"

type System struct {
	client *Client
}

type HealthResp struct {
	Status string
}

type InfoResp struct {
	NomadAddress              string
	PolicyEngine              string
	StorageBackend            string
	InternalAutoScalingEngine bool
	StrictPolicyChecking      bool
}

// LeaderResp is the response from the Leader API call.
type LeaderResp struct {
	IsSelf               bool
	HAEnabled            bool
	LeaderAddress        string
	LeaderClusterAddress string
}

func (c *Client) System() *System {
	return &System{client: c}
}

func (s *System) Health() (*HealthResp, error) {
	var resp HealthResp
	err := s.client.get("/v1/system/health", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *System) Info() (*InfoResp, error) {
	var resp InfoResp
	err := s.client.get("/v1/system/info", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *System) Metrics() (*metrics.MetricsSummary, error) {
	var resp metrics.MetricsSummary
	err := s.client.get("/v1/system/metrics", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *System) Leader() (*LeaderResp, error) {
	var resp LeaderResp
	err := s.client.get("/v1/system/leader", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
