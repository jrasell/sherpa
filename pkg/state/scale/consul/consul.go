package consul

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/jrasell/sherpa/pkg/client"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/jrasell/sherpa/pkg/state/scale"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ scale.Backend = (*StateBackend)(nil)

type StateBackend struct {
	basePath         string
	eventsPath       string
	latestEventsPath string
	gcThreshold      int64
	logger           zerolog.Logger

	kv *api.KV
}

func NewStateBackend(log zerolog.Logger, path string) scale.Backend {
	consul, _ := client.NewConsulClient()

	return &StateBackend{
		basePath:         path + "state/",
		eventsPath:       path + "state/" + "events/",
		latestEventsPath: path + "state/" + "latest-events/",
		gcThreshold:      scale.GarbageCollectionThreshold,
		logger:           log,
		kv:               consul.KV(),
	}
}

func (s StateBackend) GetScalingEvents() (map[uuid.UUID]map[string]*state.ScalingEvent, error) {
	kv, _, err := s.kv.List(s.eventsPath, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := make(map[uuid.UUID]map[string]*state.ScalingEvent)

	for i := range kv {
		keyState := &state.ScalingEvent{}

		if err := json.Unmarshal(kv[i].Value, keyState); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
		}

		keySplit := strings.Split(kv[i].Key, "/")

		id, err := uuid.FromString(keySplit[len(keySplit)-2])
		if err != nil {
			return nil, errors.Wrap(err, "failed to get UUID from string")
		}

		out[id] = map[string]*state.ScalingEvent{keySplit[len(keySplit)-1]: keyState}
	}

	return out, nil
}

func (s StateBackend) GetScalingEvent(id uuid.UUID) (map[string]*state.ScalingEvent, error) {
	kv, _, err := s.kv.List(s.eventsPath+id.String(), nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := make(map[string]*state.ScalingEvent)

	for i := range kv {
		s := &state.ScalingEvent{}
		if err := json.Unmarshal(kv[i].Value, s); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
		}

		keySplit := strings.Split(kv[i].Key, "/")
		out[keySplit[len(keySplit)-1]] = s
	}

	return out, nil
}

func (s StateBackend) PutScalingEvent(job string, event *state.ScalingEventMessage) error {

	sEntry := &state.ScalingEvent{
		EvalID:  event.EvalID,
		Source:  event.Source,
		Time:    event.Time,
		Status:  event.Status,
		Details: state.EventDetails{Count: event.Count, Direction: event.Direction},
	}

	marshal, err := json.Marshal(sEntry)
	if err != nil {
		return err
	}

	// Write the event to the general store.
	ePair := &api.KVPair{
		Key:   fmt.Sprintf("%s%s/%s:%s", s.eventsPath, event.ID.String(), job, event.GroupName),
		Value: marshal,
	}
	if _, err = s.kv.Put(ePair, nil); err != nil {
		return err
	}

	// Write the new event to the latest store.
	lePair := &api.KVPair{
		Key:   fmt.Sprintf("%s%s/%s:%s", s.latestEventsPath, event.ID.String(), job, event.GroupName),
		Value: marshal,
	}
	_, err = s.kv.Put(lePair, nil)
	return err
}

func (s StateBackend) RunGarbageCollection() {
	kv, _, err := s.kv.List(s.eventsPath, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("GC failed to list events in backend store")
	}

	if kv == nil {
		return
	}

	gc := time.Now().UTC().UnixNano() - s.gcThreshold

	for i := range kv {

		ss := &state.ScalingEvent{}

		if err := json.Unmarshal(kv[i].Value, ss); err != nil {
			s.logger.Error().Err(err).Msg("GC failed to unmarshal event for inspection")
			break
		}

		if ss.Time < gc {

			// Unlike the in-memory, we currently delete keys which have passed the expiration
			// threshold. Delete vs. re-create has not been benchmarked, but my initial opinion is
			// that delete will be more efficient and is at least easier for the MVP.
			if _, err := s.kv.Delete(kv[i].Key, nil); err != nil {
				s.logger.Error().
					Str("key", kv[i].Key).
					Err(err).
					Msg("GC failed to delete stale event in backend store")
			}
		}
	}
}
