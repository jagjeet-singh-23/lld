package eventbus

import (
	"errors"
	"testing"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

func TestEventBusPublishNotifiesAttachedObservers(t *testing.T) {
	bus := New()
	event := interfaces.Event{EventType: "order.accepted"}
	observer := &recordingObserver{}

	if err := bus.Attach(event, observer); err != nil {
		t.Fatalf("Attach() error = %v", err)
	}

	if err := bus.Publish(event); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	if len(observer.events) != 1 {
		t.Fatalf("observer received %d events, want 1", len(observer.events))
	}
	if observer.events[0].EventType != event.EventType {
		t.Fatalf("EventType = %q, want %q", observer.events[0].EventType, event.EventType)
	}
}

func TestEventBusDetachStopsNotifications(t *testing.T) {
	bus := New()
	event := interfaces.Event{EventType: "order.cancelled"}
	observer := &recordingObserver{}

	if err := bus.Attach(event, observer); err != nil {
		t.Fatalf("Attach() error = %v", err)
	}
	if err := bus.Detach(event, observer); err != nil {
		t.Fatalf("Detach() error = %v", err)
	}

	if err := bus.Publish(event); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	if len(observer.events) != 0 {
		t.Fatalf("observer received %d events, want 0", len(observer.events))
	}
}

func TestEventBusPublishReturnsObserverError(t *testing.T) {
	bus := New()
	event := interfaces.Event{EventType: "order.ready"}
	wantErr := errors.New("observer failed")

	if err := bus.Attach(event, failingObserver{err: wantErr}); err != nil {
		t.Fatalf("Attach() error = %v", err)
	}

	err := bus.Publish(event)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Publish() error = %v, want %v", err, wantErr)
	}
}

func TestEventBusPublishWithoutSubscribersIsNoop(t *testing.T) {
	bus := New()
	event := interfaces.Event{EventType: "order.delivered"}

	if err := bus.Publish(event); err != nil {
		t.Fatalf("Publish() error = %v, want nil", err)
	}
}

type recordingObserver struct {
	events []interfaces.Event
}

func (o *recordingObserver) Update(event interfaces.Event) error {
	o.events = append(o.events, event)
	return nil
}

type failingObserver struct {
	err error
}

func (o failingObserver) Update(interfaces.Event) error {
	return o.err
}
