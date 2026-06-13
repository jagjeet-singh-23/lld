package interfaces

import "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/dto"

type IOrderService interface {
	PlaceOrder(orderID string) error
	AcceptOrder(orderID string) error
	PrepareOrder(orderID string) error
	ReadyOrder(orderID string) error
	PickOrder(orderID string) error
	DeliverOrder(orderID string) error
	CancelOrder(orderID string) error
	GetOrder(orderID string) (*dto.OrderCtx, error)
}
