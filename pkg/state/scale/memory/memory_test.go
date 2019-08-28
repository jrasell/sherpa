package memory

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
	"github.com/jrasell/sherpa/pkg/state/scale"
	"github.com/stretchr/testify/assert"
)

func Test_MemoryStateBackend(t *testing.T) {
	newBackend := NewStateBackend()

	// Generate event1.
	eventTime1 := time.Now().UnixNano()
	event1 := generateTestEvent(eventTime1)

	// Write event1 to our backend state store.
	err := newBackend.PutScalingEvent("test_job_name_1", event1)
	assert.Nil(t, err)

	// Attempt to read event1 out of the state store.
	actualEvent1, err := newBackend.GetScalingEvent(event1.ID)
	expectedEvent1 := map[string]*state.ScalingEvent{
		"test_job_name_1:test_group_name": convertMessageToStateRepresentation(event1),
	}

	// Check the expected result of reading event1.
	assert.Nil(t, err)
	assert.Equal(t, expectedEvent1, actualEvent1)

	// Generate event2.
	eventTime2 := time.Now().UnixNano()
	event2 := generateTestEvent(eventTime2)

	// Write event2 to our backend state store.
	err = newBackend.PutScalingEvent("test_job_name_2", event2)
	assert.Nil(t, err)

	// Attempt to read event2 out of the state store.
	actualEvent2, err := newBackend.GetScalingEvent(event2.ID)
	expectedEvent2 := map[string]*state.ScalingEvent{
		"test_job_name_2:test_group_name": convertMessageToStateRepresentation(event2),
	}

	// Check the expected result of reading event2.
	assert.Nil(t, err)
	assert.Equal(t, expectedEvent2, actualEvent2)

	// Attempt to read out entire state out.
	actualStateRead1, err := newBackend.GetScalingEvents()

	// Check the expected results of reading back the whole state.
	expectedStateRead1 := map[uuid.UUID]map[string]*state.ScalingEvent{
		event1.ID: {"test_job_name_1:test_group_name": convertMessageToStateRepresentation(event1)},
		event2.ID: {"test_job_name_2:test_group_name": convertMessageToStateRepresentation(event2)}}

	assert.Nil(t, err)
	assert.Len(t, actualStateRead1, 2)
	assert.Equal(t, expectedStateRead1, actualStateRead1)

	// Generate an event with a time long in the past.
	eventTime3 := time.Now().UnixNano() - (scale.GarbageCollectionThreshold * 2)
	event3 := generateTestEvent(eventTime3)

	// Write event3 to our backend state store.
	err = newBackend.PutScalingEvent("test_job_name_3", event3)
	assert.Nil(t, err)

	// Attempt to read event3 out of the state store.
	actualEvent3, err := newBackend.GetScalingEvent(event3.ID)
	expectedEvent3 := map[string]*state.ScalingEvent{
		"test_job_name_3:test_group_name": convertMessageToStateRepresentation(event3),
	}

	// Check the expected result of reading event3.
	assert.Nil(t, err)
	assert.Equal(t, expectedEvent3, actualEvent3)

	// Trigger the garbage collector.
	newBackend.RunGarbageCollection()

	// Check the event has been removed from the state.
	gcEvent1, err := newBackend.GetScalingEvent(event3.ID)
	assert.Nil(t, err)
	assert.Nil(t, gcEvent1)

	// Attempt to read out entire state out.
	actualStateRead2, err := newBackend.GetScalingEvents()

	// Check the expected results of reading back the whole state.
	expectedStateRead2 := map[uuid.UUID]map[string]*state.ScalingEvent{
		event1.ID: {"test_job_name_1:test_group_name": convertMessageToStateRepresentation(event1)},
		event2.ID: {"test_job_name_2:test_group_name": convertMessageToStateRepresentation(event2)}}

	assert.Nil(t, err)
	assert.Len(t, actualStateRead2, 2)
	assert.Equal(t, expectedStateRead2, actualStateRead2)
}

func generateTestEvent(t int64) *state.ScalingEventMessage {
	id, _ := uuid.NewV4()

	return &state.ScalingEventMessage{
		ID:        id,
		GroupName: "test_group_name",
		EvalID:    id.String(),
		Source:    state.SourceAPI,
		Time:      t,
		Status:    state.StatusCompleted,
		Count:     1,
		Direction: "in",
	}
}

func convertMessageToStateRepresentation(event *state.ScalingEventMessage) *state.ScalingEvent {
	return &state.ScalingEvent{
		EvalID:  event.EvalID,
		Source:  event.Source,
		Time:    event.Time,
		Status:  event.Status,
		Details: state.EventDetails{Count: event.Count, Direction: event.Direction},
	}
}
