package aggregates

import "errors"

var (
	// ErrEventTypeNotFound is returned when an event type is not found
	ErrEventTypeNotFound = errors.New("event type not found")
)
