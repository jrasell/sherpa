package memory

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/jrasell/sherpa/pkg/state/cluster"
)

type ClusterBackend struct {
	leaderInfo  map[uuid.UUID]*state.ClusterMember
	leaderLock  *sync.RWMutex
	clusterInfo *state.ClusterInfo
	clusterLock *sync.RWMutex
}

type ClusterLock struct {
	value    string
	held     bool
	leaderCh chan struct{}
	l        *sync.Mutex
}

func NewStateBackend() cluster.Backend {
	return &ClusterBackend{
		leaderInfo: make(map[uuid.UUID]*state.ClusterMember),
	}
}

func (c ClusterBackend) DeleteLeaderEntries(uuid uuid.UUID) {
	c.leaderLock.Lock()
	defer c.leaderLock.Unlock()

	for id := range c.leaderInfo {
		if id != uuid {
			delete(c.leaderInfo, id)
		}
	}
}

func (c ClusterBackend) DeleteLeaderEntry(uuid uuid.UUID) error {
	c.leaderLock.Lock()
	defer c.leaderLock.Unlock()

	if _, ok := c.leaderInfo[uuid]; ok {
		delete(c.leaderInfo, uuid)
	}
	return nil
}

func (c ClusterBackend) GetClusterInfo() (*state.ClusterInfo, error) {
	c.clusterLock.RLock()
	info := c.clusterInfo
	c.clusterLock.RUnlock()
	return info, nil
}

func (c ClusterBackend) PutClusterInfo(info *state.ClusterInfo) error {
	c.clusterLock.Lock()
	c.clusterInfo = info
	c.clusterLock.Unlock()
	return nil
}

func (c ClusterBackend) PutClusterLeader(leader *state.ClusterMember) error {
	c.leaderLock.Lock()
	c.leaderInfo[leader.ID] = leader
	c.leaderLock.Unlock()
	return nil
}

func (c ClusterBackend) GetClusterLeader(id string) (*state.ClusterMember, error) {
	c.leaderLock.RLock()
	defer c.leaderLock.Unlock()

	uid, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	if val, ok := c.leaderInfo[uid]; ok {
		return val, nil
	}
	return nil, nil
}

func (c ClusterBackend) Lock(value string) (cluster.BackendLock, error) {
	return &ClusterLock{value: value}, nil
}

func (c ClusterBackend) SupportsHA() bool {
	return false
}

func (c ClusterLock) Acquire(stopCh <-chan struct{}) (<-chan struct{}, error) {
	c.l.Lock()
	defer c.l.Unlock()
	if c.held {
		return nil, fmt.Errorf("lock already held")
	}

	c.held = true
	c.leaderCh = make(chan struct{})
	return c.leaderCh, nil
}

func (c ClusterLock) Release() error {
	c.l.Lock()
	defer c.l.Unlock()

	if !c.held {
		return nil
	}

	close(c.leaderCh)
	c.leaderCh = nil
	c.held = false
	return nil
}

func (c ClusterLock) Value() (bool, string, error) {
	c.l.Lock()
	val := c.value
	c.l.Unlock()
	return true, val, nil
}
