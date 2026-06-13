package dto

import "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"

type OrderCtx struct {
	OrderID      string
	CurrentState enums.State
}

func NewOrderCtx(orderID string, currentState enums.State) *OrderCtx {
	return &OrderCtx{
		OrderID:      orderID,
		CurrentState: currentState,
	}
}
