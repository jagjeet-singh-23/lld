package enums

type State string

const (
	Placed    State = "placed"
	Accepted  State = "accepted"
	Preparing State = "preparing"
	Ready     State = "ready"
	PickedUp  State = "picked_up"
	Delivered State = "delivered"
	Cancelled State = "cancelled"
)
