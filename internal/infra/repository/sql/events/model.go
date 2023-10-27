package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jcleira/go-template/internal/app/aggregates"
)

// dbEvent is the database representation of an aggregates.Event.
type dbEvent struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Type      string    `db:"type"`
	Details   []byte    `db:"details"`
}

// toDomain converts a dbEvent to an aggregates.Event, for the specific event
// details it performs a JSON unmarshal from the DB stored []byte.
func (dbe dbEvent) toDomain() (aggregates.Event[any], error) {
	var details any
	switch aggregates.EventType(dbe.Type) {
	case aggregates.EventTypeCreated:
		details = aggregates.CreatedEventDetails{}
		if err := json.Unmarshal(dbe.Details, &details); err != nil {
			return aggregates.Event[any]{}, fmt.Errorf("json.Unmarshal, err: %w", err)
		}
	default:
		return aggregates.Event[any]{}, fmt.Errorf(
			"%w %s", aggregates.ErrEventTypeNotFound, dbe.Type)
	}

	return aggregates.Event[any]{
		ID:      dbe.ID,
		Type:    aggregates.EventType(dbe.Type),
		Details: details,
	}, nil
}

// dbAlertFromDomain converts an aggregates.Event to a dbEvent.
func dbAlertFromDomain(event aggregates.Event[any]) (dbEvent, error) {
	detailsBytes, err := json.Marshal(event.Details)
	if err != nil {
		return dbEvent{}, fmt.Errorf("json.Marshal, err: %w", err)
	}

	return dbEvent{
		ID:      event.ID,
		Type:    string(event.Type),
		Details: detailsBytes,
	}, nil
}

// dbEvents is a slice of dbEvent.
type dbEvents []dbEvent

// toDomain converts a slice of dbEvent to a slice of aggregates.Event.
func (dbes dbEvents) toDomain() ([]aggregates.Event[any], error) {
	events := make([]aggregates.Event[any], len(dbes))
	for i, dbe := range dbes {
		var err error
		if events[i], err = dbe.toDomain(); err != nil {
			return nil, fmt.Errorf("dbe.toDomain, err: %w", err)
		}
	}

	return events, nil
}
