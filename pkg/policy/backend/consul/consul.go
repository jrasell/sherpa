package consul

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/jrasell/sherpa/pkg/client"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ backend.PolicyBackend = (*PolicyBackend)(nil)

type PolicyBackend struct {
	path   string
	logger zerolog.Logger

	kv *api.KV
}

func NewConsulPolicyBackend(log zerolog.Logger, path string) backend.PolicyBackend {
	consul, _ := client.NewConsulClient()

	return &PolicyBackend{
		path:   path + "policies/",
		logger: log,
		kv:     consul.KV(),
	}
}

func (p *PolicyBackend) GetPolicies() (map[string]map[string]*policy.GroupScalingPolicy, error) {
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
		out[keySplit[len(keySplit)-2]] = map[string]*policy.GroupScalingPolicy{keySplit[len(keySplit)-1]: keyPolicy}
	}

	return out, nil
}

func (p *PolicyBackend) GetJobPolicy(job string) (map[string]*policy.GroupScalingPolicy, error) {
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
	kv, _, err := p.kv.Get(p.path+job+"/"+group, nil)
	if err != nil {
		return nil, err
	}

	out := &policy.GroupScalingPolicy{}

	if err := json.Unmarshal(kv.Value, out); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
	}

	return out, nil
}

func (p *PolicyBackend) PutJobPolicy(job string, groupPolicies map[string]*policy.GroupScalingPolicy) error {
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
	_, err := p.kv.DeleteTree(p.path+job, nil)
	return err
}

func (p *PolicyBackend) DeleteJobGroupPolicy(job, group string) error {
	_, err := p.kv.Delete(p.path+job+"/"+group, nil)
	return err
}
