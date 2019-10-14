package scale

import (
	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
)

// Backend is the interface required for a state storage backend. A state storage backend is used
// to durably store job scaling state outside of Sherpa.
type Backend interface {

	// GetLatestScalingEvents is used to pull all the currently stored latest scaling events from
	// the storage backend.
	GetLatestScalingEvents() (map[string]*state.ScalingEvent, error)

	// GetLatestScalingEvent is used to pull the latest scaling event for the particular job group
	// from the storage backend if we have a record.
	GetLatestScalingEvent(job, group string) (*state.ScalingEvent, error)

	// GetScalingEvents returns all scaling events held within the state.
	GetScalingEvents() (map[uuid.UUID]map[string]*state.ScalingEvent, error)

	// GetScalingEvent is used to find an individual event in the state.
	GetScalingEvent(id uuid.UUID) (map[string]*state.ScalingEvent, error)

	// PutScalingEvent is used to update the state with a new scaling event. When implementing this
	// function, care should be taken to ensure both the Events and LatestEvents fields are
	// manipulated.
	PutScalingEvent(string, *state.ScalingEventMessage) error

	// RunGarbageCollection triggers are run of the state event garbage collection which is used to
	// clear up old state entries. This ensures the state backend doesn't just continually grow.
	RunGarbageCollection()
}

const (
	// GarbageCollectionThreshold is a nano-second time, which dictates the threshold for state
	// entries to be declared stale. The current value 86400000000000 is 24 hours.
	GarbageCollectionThreshold int64 = 86400000000000
)
