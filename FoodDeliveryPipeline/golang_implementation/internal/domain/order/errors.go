package order

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
)

type DuplicateOrderError struct {
	OrderID string
}

func (e DuplicateOrderError) Error() string {
	return fmt.Sprintf(
		"unable to place duplicate order for orderID: %s",
		e.OrderID,
	)
}

type OrderNotFound struct {
	OrderID string
}

func (e OrderNotFound) Error() string {
	return fmt.Sprintf("unable to find order: %s", e.OrderID)
}

type InvalidOrderTransition struct {
	FromState enums.State
	ToState   enums.State
}

var _ error = InvalidOrderTransition{}

func (e InvalidOrderTransition) Error() string {
	return fmt.Sprintf("invalid order transition from %s to %s", e.FromState, e.ToState)
}
