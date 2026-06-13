package utils

import (
	"time"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

func NewEvent(eventType string, metadata map[string]any) interfaces.Event {
	return interfaces.Event{
		EventType: eventType,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}
