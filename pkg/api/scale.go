package api

import (
	"fmt"
	"strconv"
)

type Scale struct {
	client *Client
}

type ScaleReq struct {
	Count int
}

type ScaleResp struct {
	EvaluationID string
}

func (c *Client) Scale() *Scale {
	return &Scale{client: c}
}

func (s *Scale) JobGroupOut(job, group string, count int) (*ScaleResp, error) {
	var resp ScaleResp

	path := fmt.Sprintf("/v1/scale/out/%s/%s", job, group)

	var q QueryOptions
	if count > 0 {
		q.Params = make(map[string]string)
		q.Params["count"] = strconv.Itoa(count)
	}

	err := s.client.put(path, nil, &resp, &q)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Scale) JobGroupIn(job, group string, count int) (*ScaleResp, error) {
	var resp ScaleResp

	var q QueryOptions
	if count > 0 {
		q.Params = make(map[string]string)
		q.Params["count"] = strconv.Itoa(count)
	}

	path := fmt.Sprintf("/v1/scale/in/%s/%s", job, group)

	err := s.client.put(path, nil, &resp, &q)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
