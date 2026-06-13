package interfaces

type EventBus interface {
	Publish(event Event) error
	Attach(event Event, observer Observer) error
	Detach(event Event, observer Observer) error
}
