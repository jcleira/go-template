package events

import (
	"time"
)

type dbEvent struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Type      string    `db:"type"`
	Details   string    `db:"value"`
}

func (dbe dbEvent) ToDomain() aggregates.Alert {
	switch dbe.Type {
	case aggregates.EventTypeCreated:
		details := aggregates.CreatedEventDetails{}
		if err := json.Unmarshal(dbe.Details, &details); err != nil {
			return event, fmt.Errorf("json.Unmarshal, err: %w", err)
		}
	}

	return aggregates.Alert{
		ID:          dbe.ID,
		Type:        dbe.Type,
		Details:     details
	}
}

func dbAlertFromDomain(alert aggregates.Alert) dbEvent {
	return dbEvent{
		ID:          alert.ID,
		TriggeredAt: alert.TriggeredAt,
		Type:        alert.Type,
		Details:     alert.Details,
	}
}

type dbAlerts []dbEvent

func (dbe dbAlerts) ToDomain() ([]aggregates.Alert, error) {
	alerts := make([]aggregates.Alert, len(dbe))
	for i, alert := range dbe {
		alerts[i] = alert.ToDomain()
	}

	return alerts, nil
}
