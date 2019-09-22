package cluster

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state/cluster"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Member struct { // nolint:maligned

	// local member information.
	id      uuid.UUID
	addr    string
	advAddr string

	// clusterStorage interface is used to write and read updates from the backend.
	clusterStorage   cluster.Backend
	clusterStorageHA bool

	clusterLock cluster.BackendLock

	stateLock sync.RWMutex
	logger    zerolog.Logger

	// clusterName
	clusterName string

	// standby is used to determine whether the Sherpa server in question is currently in standby
	// mode or not. Standby is an instance which does not currently hold the leadership lock.
	standby bool

	// clusterLeader information is protected by a mutex so it can be safely read and written. The
	// information here is periodically checked to ensure it is correct.
	clusterLeaderLock          sync.RWMutex
	clusterLeaderID            uuid.UUID
	clusterLeaderAddr          string
	clusterLeaderAdvertiseAddr string

	// stopChan is used by the cluster member to coordinate the stopping of background tasks.
	stopChan chan struct{}

	// UpdateChan is used to publish leadership updates to the server. This allows the server to
	// coordinate tasks such as scaling state garbage collection and the autoscaler, both of which
	// should only be run by the leader.
	UpdateChan chan *MembershipUpdate
}

// MembershipUpdate is the information used to update the server about leadership changes.
type MembershipUpdate struct {

	// IsLeader is used to determine whether or not the Sherpa server is acting as the cluster
	// leader.
	IsLeader bool

	// Msg is an optional message to pass which can be useful when logging update events.
	Msg string
}

func NewMember(log zerolog.Logger, store cluster.Backend, addr, advAddr, name string) (*Member, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup new member ID")
	}

	m := Member{
		clusterStorage:   store,
		clusterStorageHA: store.SupportsHA(),
		logger:           log,
		id:               id,
		addr:             addr,
		advAddr:          advAddr,
		clusterName:      name,
		UpdateChan:       make(chan *MembershipUpdate),
		stopChan:         make(chan struct{}),
		standby:          true,
	}

	if err := m.setupCluster(); err != nil {
		return nil, errors.Wrap(err, "failed to setup cluster member")
	}

	// Set the cluster member logger to contain contextual information including the name and ID.
	m.logger = log.With().
		Str("cluster-member-id", m.id.String()).
		Str("cluster-name", m.clusterName).Logger()

	return &m, nil
}

func (m *Member) IsHA() bool { return m.clusterStorageHA }
