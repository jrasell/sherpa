package state

import (
	"github.com/gofrs/uuid"
)

// ScalingState is the internal Sherpa scaling state used to track scaling events invoked.
type ScalingState struct {
	// Events stores all scaling events to have been invoked through the Sherpa server. These
	// events can be from any source. Each scaling event will be registered with a UUID; each event
	// can encompass a single job, and any number of groups within that job.
	Events map[uuid.UUID]map[string]*ScalingEvent

	// LatestEvents holds the most recent scaling event to occur for each job group. The map key
	// takes the form of job-name:group-name. When a new scaling action takes place for a
	// particular job group, the entry here should be overwritten. This provides and fast way to
	// lookup last events and is currently ignored from GC.
	LatestEvents map[string]*ScalingEvent
}

// ScalingEvent represents a single scaling event state entry that is persisted to the backend
// store.
type ScalingEvent struct {
	// ID is the scaling ID.
	ID uuid.UUID

	// EvalID is the Nomad evaluation ID which was created as a result of submitting the updated
	// job to the Nomad API.
	EvalID string

	// Source shows the origin source of the scaling event.
	Source Source

	// Time is a UnixNano timestamp declaring when the scaling event trigger took place.
	Time int64

	// Status is the end status of the scaling event.
	Status Status

	// Details contains information about exactly what action was taken on the job group during the
	// scaling event.
	Details EventDetails

	Meta map[string]string
}

// EventDetails contains information to describe what changes took place during the scaling action.
type EventDetails struct {
	// Count is the number by which the group was changed.
	Count int

	// Direction is direction in which the scaling took place. This can be either in or out
	// representing subtraction and addition respectively.
	Direction string
}

// ScalingEventMessage is the message sent to the state writer containing all the required
// information to construct the persistent state entry.
type ScalingEventMessage struct {
	ID        uuid.UUID
	GroupName string
	EvalID    string
	Source    Source
	Time      int64
	Status    Status
	Count     int
	Direction string
	Meta      map[string]string
}

// Source represents how the scaling action was invoked.
type Source string

const (
	// SourceAPI is a scaling event invocation via the HTTP API.
	SourceAPI Source = "API"

	// SourceInternalAutoscaler is a scaling event invoked by the internal autoscaler.
	SourceInternalAutoscaler Source = "InternalAutoscaler"
)

func (s Source) String() string { return string(s) }

// Status represents whether the scaling event was classed as successful or not. Currently this is
// dependant on if the job managed to be submitted to the Nomad API as Sherpa does not perform any
// checking after.
type Status string

const (
	// StatusCompleted means the changes to the Nomad job specification due to a scaling
	// requirement were registered to the Nomad API without error.
	StatusCompleted = "Completed"

	// StatusFailed means there was an error calling the Nomad API when attempting to register
	// the job which contained altered groups as a result of a scaling event.
	StatusFailed = "Failed"
)

func (s Status) String() string { return string(s) }
