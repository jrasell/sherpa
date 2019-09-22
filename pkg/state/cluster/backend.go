package cluster

import (
	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
)

// Backend is the interface used to govern most leadership state actives, apart from the actual
// synchronisation which is handled by the BackendLock interface.
type Backend interface {

	// DeleteLeaderEntries is used to prune old leader entries within the backend store. The passed
	// ID is the exception and should be left intact if it is found to be present.
	DeleteLeaderEntries(uuid uuid.UUID)

	// DeleteLeaderEntry will delete the entry of the passed ID in the backend store if it exists.
	DeleteLeaderEntry(uuid uuid.UUID) error

	// GetClusterInfo will attempt to return the stored cluster information from the backend. If
	// there is no information then nil should be returned indicating the Sherpa server is part of
	// a new cluster.
	GetClusterInfo() (*state.ClusterInfo, error)

	// PutClusterInfo can be used to upsert ClusterInfo data into the backend.
	PutClusterInfo(info *state.ClusterInfo) error

	// PutClusterLeader is used to write the information regarding the current cluster leader to
	// the backend store. This contains information relevant for request redirection and such.
	PutClusterLeader(leader *state.ClusterMember) error

	// GetClusterLeader attempts to read the current cluster leader information based on the passed
	// member ID.
	GetClusterLeader(id string) (*state.ClusterMember, error)

	// Lock is used for mutual exclusion based on the passed value.
	Lock(value string) (BackendLock, error)

	// SupportsHA is used to determine whether the backend storage system supports high
	// availability features.
	SupportsHA() bool
}

// BackendLock is the locking interface which is used to perform and monitor leadership activities.
type BackendLock interface {

	// Acquire is used to acquire the given lock. The stopCh should interrupt the lock acquisition
	// attempt. The return struct is closed when leadership is lost, which can be used by clients
	// as an indication to restart the process.
	Acquire(stopCh <-chan struct{}) (<-chan struct{}, error)

	// Release is used to release the current leadership lock.
	Release() error

	// Value is used to return the value of the lock, if it is currently being held.
	Value() (bool, string, error)
}
