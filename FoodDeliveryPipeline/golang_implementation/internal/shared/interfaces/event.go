package interfaces

import "time"

type Event struct {
	EventType string
	Metadata  map[string]any
	Timestamp time.Time
}

func NewEvent(eventType string, metadata map[string]any) Event {
	return Event{
		EventType: eventType,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}
