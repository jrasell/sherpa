package consul

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/jrasell/sherpa/pkg/state/cluster"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	sessionLockName   = "Sherpa Lock"
	clusterInfoPath   = "cluster/info"
	clusterLockPath   = "cluster/lock"
	clusterLeaderPath = "cluster/leader/"
)

type ClusterBackend struct {
	client *api.Client
	kv     *api.KV
	logger zerolog.Logger

	clusterInfoPath   string
	clusterLockPath   string
	clusterLeaderPath string

	sessionTTL   string
	lockWaitTime time.Duration
}

type ClusterLock struct {
	client *api.Client
	key    string
	lock   *api.Lock
}

func NewStateBackend(log zerolog.Logger, path string, client *api.Client) cluster.Backend {
	return &ClusterBackend{
		client:            client,
		kv:                client.KV(),
		clusterInfoPath:   path + clusterInfoPath,
		clusterLockPath:   path + clusterLockPath,
		clusterLeaderPath: path + clusterLeaderPath,
		logger:            log,
		sessionTTL:        api.DefaultLockSessionTTL,
		lockWaitTime:      api.DefaultLockWaitTime,
	}
}

func (c ClusterBackend) DeleteLeaderEntries(uuid uuid.UUID) {
	keys, _, err := c.kv.Keys(c.clusterLeaderPath, "/", nil)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to list leader entries")
		return
	}

	for _, val := range keys {
		id := strings.TrimPrefix(val, c.clusterLeaderPath)
		if id != uuid.String() {
			_, err := c.kv.Delete(val, nil)
			if err != nil {
				c.logger.Error().Err(err).Msg("failed to delete leadership entry")
			}
		}
	}
}

func (c ClusterBackend) DeleteLeaderEntry(uuid uuid.UUID) error {
	key := c.clusterLeaderPath + uuid.String()
	_, err := c.kv.Delete(key, nil)
	return err
}

func (c ClusterBackend) GetClusterInfo() (*state.ClusterInfo, error) {
	kv, _, err := c.kv.Get(c.clusterInfoPath, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	info := &state.ClusterInfo{}
	if err := json.Unmarshal(kv.Value, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (c ClusterBackend) PutClusterInfo(info *state.ClusterInfo) error {
	bytes, err := json.Marshal(info)
	if err != nil {
		return err
	}

	_, err = c.kv.Put(&api.KVPair{Key: c.clusterInfoPath, Value: bytes}, nil)
	return err
}

func (c ClusterBackend) GetClusterLeader(id string) (*state.ClusterMember, error) {
	kv, _, err := c.kv.Get(c.clusterLeaderPath+id, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	mem := &state.ClusterMember{}
	if err := json.Unmarshal(kv.Value, mem); err != nil {
		return nil, err
	}
	return mem, nil
}

func (c ClusterBackend) PutClusterLeader(leader *state.ClusterMember) error {
	bytes, err := json.Marshal(leader)
	if err != nil {
		return err
	}

	kv := api.KVPair{
		Key:   c.clusterLeaderPath + leader.ID.String(),
		Value: bytes,
	}

	_, err = c.kv.Put(&kv, nil)
	return err
}

func (c ClusterBackend) Lock(value string) (cluster.BackendLock, error) {
	opts := &api.LockOptions{
		Key:            c.clusterLockPath,
		Value:          []byte(value),
		SessionName:    sessionLockName,
		MonitorRetries: 5,
		SessionTTL:     c.sessionTTL,
		LockWaitTime:   c.lockWaitTime,
	}
	lock, err := c.client.LockOpts(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lock")
	}
	cl := &ClusterLock{
		client: c.client,
		key:    c.clusterLockPath,
		lock:   lock,
	}
	return cl, nil
}

func (c ClusterBackend) SupportsHA() bool {
	return true
}

func (c ClusterLock) Acquire(stopCh <-chan struct{}) (<-chan struct{}, error) {
	return c.lock.Lock(stopCh)
}

func (c ClusterLock) Release() error {
	return c.lock.Unlock()
}

func (c ClusterLock) Value() (bool, string, error) {
	kv := c.client.KV()

	pair, _, err := kv.Get(c.key, nil)
	if err != nil {
		return false, "", err
	}
	if pair == nil {
		return false, "", nil
	}
	held := pair.Session != ""
	value := string(pair.Value)
	return held, value, nil
}
