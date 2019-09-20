package cluster

import (
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/jrasell/sherpa/pkg/state/cluster"
	"github.com/oklog/run"
)

const (

	// leaderCheckInterval the time interval by which a standby Sherpa server will check for a new
	// leader.
	leaderCheckInterval = 2500 * time.Millisecond

	// lockRetryInterval is the interval we re-attempt to acquire the HA lock if an error is
	// encountered.
	lockRetryInterval = 10 * time.Second
)

const (
	updateMsgObtainedLeadership = "obtained leadership"
	updateMsgLostLeadership     = "lost leadership"
)

func (m *Member) RunLeadershipLoop() {

	// Group collects functions and runs them concurrently. When one function returns, all funcs
	// are interrupted. This allows us to run a number of leadership tasks with consistent
	// synchronisation on closure.
	var g run.Group

	{
		// This will cause all the other actors to close when the stop channel is signaled.
		g.Add(func() error {
			<-m.stopChan
			return nil
		}, func(error) {})
	}

	{
		// Monitor for new leadership
		checkLeaderStop := make(chan struct{})

		g.Add(func() error {
			m.leaderRefresh(checkLeaderStop)
			return nil
		}, func(error) {
			close(checkLeaderStop)
			m.logger.Debug().Msg("shutting down periodic leader refresh")
		})
	}
	{
		// Wait for leadership
		leaderStopCh := make(chan struct{})

		g.Add(func() error {
			m.waitForLeadership(leaderStopCh)
			return nil
		}, func(error) {
			close(leaderStopCh)
			m.logger.Debug().Msg("shutting down leader elections")
		})
	}

	if err := g.Run(); err != nil {
		m.logger.Error().Err(err).Msg("failed to correctly start leadership loop actors")
	}
}

// Leader is used to identify the clusters currently recognised leader and identify the associated
// HTTP addresses.
func (m *Member) Leader() (bool, string, string, error) {

	// Lock our state for reading. This ensures items do not change while reading them and means
	// the data we provide is correct at the time.
	m.stateLock.RLock()

	// If we are not the standby, then return our stored information to the caller which will
	// identify this instance of Sherpa to be running as the leader.
	if !m.standby {
		m.stateLock.RUnlock()
		return true, m.addr, m.advAddr, nil
	}

	lock, err := m.clusterStorage.Lock("read")
	if err != nil {
		m.stateLock.RUnlock()
		return false, "", "", err
	}

	held, id, err := lock.Value()
	if err != nil {
		m.stateLock.RUnlock()
		return false, "", "", err
	}
	if !held {
		m.stateLock.RUnlock()
		return false, "", "", nil
	}

	m.clusterLeaderLock.RLock()
	localLeaderUUID := m.clusterLeaderID
	localLeaderAddr := m.clusterLeaderAddr
	localLeaderAdvAddr := m.clusterLeaderAdvertiseAddr
	m.clusterLeaderLock.RUnlock()

	if id == localLeaderUUID.String() && localLeaderAddr != "" && localLeaderAdvAddr != "" {
		m.stateLock.RUnlock()
		return false, localLeaderAddr, localLeaderAdvAddr, nil
	}
	m.logger.Debug().Msg("found new leadership information, updating internal references")

	defer m.stateLock.RUnlock()
	m.clusterLeaderLock.Lock()
	defer m.clusterLeaderLock.Unlock()

	// Validate base conditions again
	if id == m.clusterLeaderID.String() && m.clusterLeaderAddr != "" && localLeaderAdvAddr != "" {
		return false, localLeaderAddr, localLeaderAdvAddr, nil
	}

	leader, err := m.clusterStorage.GetClusterLeader(id)
	if err != nil {
		return false, "", "", err
	}

	if leader == nil {
		return false, "", "", nil
	}

	m.clusterLeaderID = leader.ID
	m.clusterLeaderAddr = leader.Addr
	m.clusterLeaderAdvertiseAddr = leader.AdvertiseAddr

	return false, leader.Addr, leader.AdvertiseAddr, nil
}

// ClearLeadership is used to coordinate the shutdown of the cluster membership processes so we can
// exist cleanly. This allows other servers to quickly take over as leader and therefore resume
// cluster operations.
func (m *Member) ClearLeadership() {
	m.logger.Info().Msg("shutting down leadership handler")
	m.stopChan <- struct{}{}
	close(m.stopChan)

	if m.clusterLock != nil {
		if err := m.clusterLock.Release(); err != nil {
			m.logger.Error().Err(err).Msg("failed to gracefully release leadership lock")
		}
	}

	if err := m.removeAsLeader(m.id); err != nil {
		m.logger.Error().Err(err).Msg("failed to gracefully remove leadership state entry")
	}
}

func (m *Member) acquireLock(lock cluster.BackendLock, stopCh <-chan struct{}) <-chan struct{} {
	for {
		// Attempt lock acquisition.
		leaderLostCh, err := lock.Acquire(stopCh)
		if err == nil {
			return leaderLostCh
		}

		// Retry the acquisition.
		m.logger.Error().Err(err).Msg("failed to acquire lock")
		select {
		case <-time.After(lockRetryInterval):
		case <-stopCh:
			return nil
		}
	}
}

func (m *Member) waitForLeadership(stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			m.logger.Debug().Msg("stop channel triggered")
			return
		default:
		}

		lock, err := m.clusterStorage.Lock(m.id.String())
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to get lock on storage backend")
			return
		}

		// Attempt the acquisition
		leaderLostCh := m.acquireLock(lock, stopCh)
		if leaderLostCh == nil {
			m.logger.Debug().Msg("failed to acquire leadership lock")
			return
		}

		if stopped := grabLockOrStop(m.stateLock.Lock, m.stateLock.Unlock, stopCh); stopped {
			if err := lock.Release(); err != nil {
				m.logger.Error().Err(err).Msg("failed to release the held leadership lock")
			}
			return
		}

		// Store the lock so that we can manually clear it later if needed
		m.clusterLock = lock

		//
		if err := m.setAsLeader(); err != nil {
			m.clusterLock = nil
			if err := lock.Release(); err != nil {
				m.logger.Error().Err(err).Msg("failed to release the held leadership lock")
			}
			m.stateLock.Unlock()
			m.logger.Error().Err(err).Msg("failed to set server as leader")
			continue
		}
		m.logger.Info().Msg("server is now acting as Sherpa cluster leader")

		// At this point we are now acting as cluster leader. Inform the server and set our state
		// to declare we are no longer a standby.
		m.standby = false
		m.UpdateChan <- &MembershipUpdate{IsLeader: true, Msg: updateMsgObtainedLeadership}
		m.stateLock.Unlock()

		// Block on either being stopped, or in the event that we lose leadership.
		select {
		case <-leaderLostCh:
			// If we have lost leadership, inform the server so that Sherpa process can be stopped,
			// then continue through the rest of the loop.
			m.logger.Warn().Msg("cluster leadership has been lost")
			m.UpdateChan <- &MembershipUpdate{IsLeader: false, Msg: updateMsgLostLeadership}

		case <-stopCh:
			// If we are told to stop, then we should just return here. Another process is
			// responsible for performing shutdown cleanup.
			return
		}

		{
			// Grab lock if we are not stopped.
			stopped := grabLockOrStop(m.stateLock.Lock, m.stateLock.Unlock, stopCh)

			// Mark that we are now standby.
			m.standby = true

			if err := m.removeAsLeader(m.id); err != nil {
				m.logger.Error().Err(err).Msg("clearing leader advertisement failed")
			}

			if err := m.clusterLock.Release(); err != nil {
				m.logger.Err(err).Msg("failed to release cluster lock")
			}
			m.clusterLock = nil

			// If we are stopped return, otherwise unlock the statelock and restart the leadership
			// loop.
			if stopped {
				return
			}
			m.stateLock.Unlock()
		}
	}
}

func grabLockOrStop(lockFunc, unlockFunc func(), stopCh chan struct{}) (stopped bool) {

	// Grab the lock as we need it for cluster setup, which needs to happen
	// before advertising;
	lockGrabbedCh := make(chan struct{})
	go func() {
		// Grab the lock
		lockFunc()
		// If stopChan has been closed, which only happens while the stateLock is held, we have
		// actually terminated, so we just instantly give up the lock, otherwise we notify that
		// it's ready for consumption.
		select {
		case <-stopCh:
			unlockFunc()
		default:
			close(lockGrabbedCh)
		}
	}()

	select {
	case <-stopCh:
		return true
	case <-lockGrabbedCh:
		// We now have the lock and can use it
	}

	return false
}

func (m *Member) leaderRefresh(stopCh chan struct{}) {
	opCount := new(int32)
	for {
		select {
		case <-time.After(leaderCheckInterval):
			count := atomic.AddInt32(opCount, 1)
			if count > 1 {
				atomic.AddInt32(opCount, -1)
				continue
			}
			go func() {
				lopCount := opCount
				_, _, _, _ = m.Leader()
				atomic.AddInt32(lopCount, -1)
			}()
		case <-stopCh:
			return
		}
	}
}

// setAsLeader is used to set the current Sherpa instance as the leader by updating the backend
// storage to reflect the leader information.
func (m *Member) setAsLeader() error {
	// Delete old cluster leader entries, ensuring we preserve our own entry. This is a maintenance
	// task and make sense to run at this point.
	go m.clusterStorage.DeleteLeaderEntries(m.id)

	leaderEntry := state.ClusterMember{ID: m.id, Addr: m.addr, AdvertiseAddr: m.advAddr}
	return m.clusterStorage.PutClusterLeader(&leaderEntry)
}

// removeAsLeader removes our backend storage entry within the leader partition.
func (m *Member) removeAsLeader(uuid uuid.UUID) error {
	return m.clusterStorage.DeleteLeaderEntry(uuid)
}
