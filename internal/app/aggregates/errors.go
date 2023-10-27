package aggregates

import "errors"

// ErrEventTypeNotFound is returned when an event type is not found.
var ErrEventTypeNotFound = errors.New("event type not found")
