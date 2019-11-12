package consul

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ backend.PolicyBackend = (*PolicyBackend)(nil)

const (
	baseKVPath = "policies/"
)

// Define our metric keys.
var (
	metricKeyGetPolicies          = []string{"policy", "consul", "get_policies"}
	metricKeyGetJobPolicy         = []string{"policy", "consul", "get_job_policy"}
	metricKeyGetJobGroupPolicy    = []string{"policy", "consul", "get_job_group_policy"}
	metricKeyPutJobPolicy         = []string{"policy", "consul", "put_job_policy"}
	metricKeyPutJobGroupPolicy    = []string{"policy", "consul", "put_job_group_policy"}
	metricKeyDeleteJobPolicy      = []string{"policy", "consul", "delete_job_policy"}
	metricKeyDeleteJobGroupPolicy = []string{"policy", "consul", "delete_job_group_policy"}
)

type PolicyBackend struct {
	path   string
	logger zerolog.Logger

	kv *api.KV
}

func NewConsulPolicyBackend(log zerolog.Logger, path string, client *api.Client) backend.PolicyBackend {
	return &PolicyBackend{
		path:   path + baseKVPath,
		logger: log,
		kv:     client.KV(),
	}
}

func (p *PolicyBackend) GetPolicies() (map[string]map[string]*policy.GroupScalingPolicy, error) {
	defer metrics.MeasureSince(metricKeyGetPolicies, time.Now())

	kv, _, err := p.kv.List(p.path, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := make(map[string]map[string]*policy.GroupScalingPolicy)

	for i := range kv {
		keyPolicy := &policy.GroupScalingPolicy{}

		if err := json.Unmarshal(kv[i].Value, keyPolicy); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
		}

		keySplit := strings.Split(kv[i].Key, "/")
		jobName := keySplit[len(keySplit)-2]
		groupName := keySplit[len(keySplit)-1]

		if _, ok := out[jobName]; !ok {
			out[jobName] = map[string]*policy.GroupScalingPolicy{}
		}

		out[jobName][groupName] = keyPolicy
	}

	return out, nil
}

func (p *PolicyBackend) GetJobPolicy(job string) (map[string]*policy.GroupScalingPolicy, error) {
	defer metrics.MeasureSince(metricKeyGetJobPolicy, time.Now())

	kv, _, err := p.kv.List(p.path+job, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := make(map[string]*policy.GroupScalingPolicy)

	for i := range kv {

		keyPolicy := &policy.GroupScalingPolicy{}

		if err := json.Unmarshal(kv[i].Value, keyPolicy); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
		}

		keySplit := strings.Split(kv[i].Key, "/")
		out[keySplit[len(keySplit)-1]] = keyPolicy
	}

	return out, nil
}

func (p *PolicyBackend) GetJobGroupPolicy(job, group string) (*policy.GroupScalingPolicy, error) {
	defer metrics.MeasureSince(metricKeyGetJobGroupPolicy, time.Now())

	kv, _, err := p.kv.Get(p.path+job+"/"+group, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := &policy.GroupScalingPolicy{}

	if err := json.Unmarshal(kv.Value, out); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
	}

	return out, nil
}

func (p *PolicyBackend) PutJobPolicy(job string, groupPolicies map[string]*policy.GroupScalingPolicy) error {
	defer metrics.MeasureSince(metricKeyPutJobPolicy, time.Now())

	var kvOpts []*api.KVTxnOp // nolint:prealloc

	for group, pol := range groupPolicies {

		marshal, err := json.Marshal(pol)
		if err != nil {
			return err
		}

		kvOpt := &api.KVTxnOp{
			Verb:  api.KVSet,
			Key:   p.path + job + "/" + group,
			Value: marshal,
		}
		kvOpts = append(kvOpts, kvOpt)
	}

	success, _, _, err := p.kv.Txn(kvOpts, nil)
	if err != nil {
		return err
	}

	if !success {
		return errors.New("failed to write job policy Consul transaction")
	}

	return nil
}

func (p *PolicyBackend) PutJobGroupPolicy(job, group string, pol *policy.GroupScalingPolicy) error {
	defer metrics.MeasureSince(metricKeyPutJobGroupPolicy, time.Now())

	marshal, err := json.Marshal(pol)
	if err != nil {
		return err
	}

	pair := &api.KVPair{
		Key:   p.path + job + "/" + group,
		Value: marshal,
	}

	_, err = p.kv.Put(pair, nil)
	return err
}

func (p *PolicyBackend) DeleteJobPolicy(job string) error {
	defer metrics.MeasureSince(metricKeyDeleteJobPolicy, time.Now())

	_, err := p.kv.DeleteTree(p.path+job, nil)
	return err
}

func (p *PolicyBackend) DeleteJobGroupPolicy(job, group string) error {
	defer metrics.MeasureSince(metricKeyDeleteJobGroupPolicy, time.Now())

	_, err := p.kv.Delete(p.path+job+"/"+group, nil)
	return err
}
