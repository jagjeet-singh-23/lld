package delivery

import "fmt"

type DeliveryError interface {
	Error() string
}

type NoDeliveryPartnerFound struct {
	OrderID string
}

func (e NoDeliveryPartnerFound) Error() string {
	return fmt.Sprintf("unable to assign a delivery partner for order: %s", e.OrderID)
}

type OrderNotReadyYet struct {
	OrderID string
}

func (e OrderNotReadyYet) Error() string {
	return fmt.Sprintf("unable to assign delivery for order: %s, order not ready yet...", e.OrderID)
}
