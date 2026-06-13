package eventbus

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

type EventBus struct {
	observers map[string][]interfaces.Observer
}

func New() *EventBus {
	return &EventBus{
		observers: make(map[string][]interfaces.Observer),
	}
}

func (b *EventBus) Attach(event interfaces.Event, observer interfaces.Observer) error {
	eventType := event.EventType
	b.observers[eventType] = append(b.observers[eventType], observer)
	return nil
}

func (b *EventBus) Detach(event interfaces.Event, observer interfaces.Observer) error {
	subscribers := b.observers[event.EventType]
	for idx, subscriber := range subscribers {
		if subscriber == observer {
			b.observers[event.EventType] = append(subscribers[:idx], subscribers[idx+1:]...)
			return nil
		}
	}

	return nil
}

func (b *EventBus) Publish(event interfaces.Event) error {
	for idx, observer := range b.observers[event.EventType] {
		if err := observer.Update(event); err != nil {
			return fmt.Errorf("notify observer %d for event %q: %w", idx, event.EventType, err)
		}
	}

	return nil
}

var _ interfaces.EventBus = (*EventBus)(nil)
