package aggregates

type EventType string

const EventTypeCreated EventType = "created"

type CreatedEventDetails struct {
	CreatedAt time.Time
}

type Event[Details any] struct {
	ID      string
	Type    EventType
	Details Details
}
