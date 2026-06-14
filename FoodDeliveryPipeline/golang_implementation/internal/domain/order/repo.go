package order

import (
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/dto"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
)

type OrdersRepo struct {
	Orders map[string]*dto.OrderCtx
}

func NewOrderRepository() *OrdersRepo {
	return &OrdersRepo{
		Orders: make(map[string]*dto.OrderCtx),
	}
}

func (r *OrdersRepo) PlaceOrder(ctx *dto.OrderCtx) error {
	if _, ok := r.Orders[ctx.OrderID]; ok {
		return DuplicateOrderError{
			OrderID: ctx.OrderID,
		}
	}
	r.Orders[ctx.OrderID] = ctx
	return nil
}

func (r *OrdersRepo) GetOrder(orderID string) (*dto.OrderCtx, error) {
	if _, ok := r.Orders[orderID]; !ok {
		return nil, OrderNotFound{
			OrderID: orderID,
		}
	}

	return r.Orders[orderID], nil
}

func (r *OrdersRepo) UpdateOrder(ctx *dto.OrderCtx) error {
	if _, ok := r.Orders[ctx.OrderID]; !ok {
		return OrderNotFound{
			OrderID: ctx.OrderID,
		}
	}

	r.Orders[ctx.OrderID] = ctx
	return nil
}

func (r *OrdersRepo) TransitionTo(orderID string, state enums.State) error {
	ctx, err := r.GetOrder(orderID)
	if err != nil {
		return err
	}

	ctx.CurrentState = state
	return r.UpdateOrder(ctx)
}

func (r *OrdersRepo) UpdateOrderState(orderID string, state enums.State) error {
	return r.TransitionTo(orderID, state)
}
