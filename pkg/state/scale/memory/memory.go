package memory

import (
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/jrasell/sherpa/pkg/state/scale"
)

var _ scale.Backend = (*StateBackend)(nil)

// Define our metric keys.
var (
	metricKeyGetEvents       = []string{"scale", "state", "memory", "get_events"}
	metricKeyGetEvent        = []string{"scale", "state", "memory", "get_event"}
	metricKeyGetLatestEvents = []string{"scale", "state", "memory", "get_latest_events"}
	metricKeyGetLatestEvent  = []string{"scale", "state", "memory", "get_latest_event"}
	metricKeyPutEvent        = []string{"scale", "state", "memory", "put_event"}
	metricKeyGC              = []string{"scale", "state", "memory", "gc"}
)

type StateBackend struct {
	gcThreshold int64
	state       *state.ScalingState
	sync.RWMutex
}

func NewStateBackend() scale.Backend {
	return &StateBackend{
		gcThreshold: scale.GarbageCollectionThreshold,
		state: &state.ScalingState{
			Events:       make(map[uuid.UUID]map[string]*state.ScalingEvent),
			LatestEvents: make(map[string]*state.ScalingEvent),
		},
	}
}

func (s *StateBackend) GetLatestScalingEvents() (map[string]*state.ScalingEvent, error) {
	defer metrics.MeasureSince(metricKeyGetLatestEvents, time.Now())

	s.RLock()
	latest := s.state.LatestEvents
	s.RUnlock()
	return latest, nil
}

func (s *StateBackend) GetLatestScalingEvent(job, group string) (*state.ScalingEvent, error) {
	defer metrics.MeasureSince(metricKeyGetLatestEvent, time.Now())

	s.RLock()
	latest := s.state.LatestEvents[job+":"+group]
	s.RUnlock()
	return latest, nil
}

func (s *StateBackend) GetScalingEvents() (map[uuid.UUID]map[string]*state.ScalingEvent, error) {
	defer metrics.MeasureSince(metricKeyGetEvents, time.Now())

	s.RLock()
	events := s.state.Events
	s.RUnlock()
	return events, nil
}

func (s *StateBackend) PutScalingEvent(job string, event *state.ScalingEventMessage) error {
	defer metrics.MeasureSince(metricKeyPutEvent, time.Now())

	s.Lock()
	defer s.Unlock()

	k := job + ":" + event.GroupName

	sEntry := &state.ScalingEvent{
		EvalID:  event.EvalID,
		Source:  event.Source,
		Time:    event.Time,
		Status:  event.Status,
		Details: state.EventDetails{Count: event.Count, Direction: event.Direction},
		Meta:    event.Meta,
	}

	s.state.Events[event.ID] = make(map[string]*state.ScalingEvent)
	s.state.Events[event.ID][k] = sEntry
	s.state.LatestEvents[k] = sEntry

	return nil
}

func (s *StateBackend) GetScalingEvent(id uuid.UUID) (map[string]*state.ScalingEvent, error) {
	defer metrics.MeasureSince(metricKeyGetEvent, time.Now())

	s.RLock()
	e := s.state.Events[id]
	s.RUnlock()
	return e, nil
}

func (s *StateBackend) RunGarbageCollection() {
	t := time.Now()
	defer metrics.MeasureSince(metricKeyGC, t)

	gc := t.UTC().UnixNano() - s.gcThreshold

	newEventState := make(map[uuid.UUID]map[string]*state.ScalingEvent)

	// Perform a read lock while performing the calculation work.
	s.RLock()

	// Iterate the event state. We do not perform GC on the latest tracked events so that we can
	// always use these in the future.
	for id, jgEvent := range s.state.Events {
		for name, event := range jgEvent {
			if event.Time > gc {
				newEventState[id] = make(map[string]*state.ScalingEvent)
				newEventState[id][name] = event
			}
		}
	}

	// Remove the read lock, and lock for writing.
	s.RUnlock()
	s.Lock()

	// Replace the internal events state with the newly built state.
	s.state.Events = newEventState
	s.Unlock()
}
