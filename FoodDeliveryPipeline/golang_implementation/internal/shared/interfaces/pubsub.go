package interfaces

type Observer interface {
	Update(event Event) error
}

type Publisher interface {
	Publish(event Event) error
}
