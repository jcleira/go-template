package aggregates

import "time"

// EvenType is a custom type to create an enumeration of type's of events.
type EventType string

// EventTypeCreated is the default/first event while creating other aggregate.
const EventTypeCreated EventType = "created"

// CreatedEventDetails define the details attributes of the created event.
type CreatedEventDetails struct {
	CreatedAt time.Time
}

// Event is the main aggregate it will be used to create specific events
// defined by the event type
type Event[Details any] struct {
	ID      string
	Type    EventType
	Details Details
}
