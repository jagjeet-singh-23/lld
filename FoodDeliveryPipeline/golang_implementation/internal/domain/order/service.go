package order

import (
	"fmt"
	"slices"
	"time"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/dto"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/utils"
)

type OrderService struct {
	repo      *OrdersRepo
	bus       interfaces.EventBus
	auditRepo *TransitionAuditRepo
}

type orderAction string

const (
	acceptOrder  orderAction = "accept"
	prepareOrder orderAction = "prepare"
	readyOrder   orderAction = "ready"
	pickOrder    orderAction = "pick"
	deliverOrder orderAction = "deliver"
	cancelOrder  orderAction = "cancel"
)

func NewOrderService(repo *OrdersRepo, bus interfaces.EventBus) *OrderService {
	return &OrderService{
		repo:      repo,
		bus:       bus,
		auditRepo: NewTransitionAuditRepo(),
	}
}

func (s *OrderService) PlaceOrder(orderID string) error {
	ctx := dto.NewOrderCtx(orderID, enums.Placed)
	err := s.repo.PlaceOrder(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Order placed successfully: %s", ctx.OrderID)
	event := utils.NewEvent(string(enums.Placed), map[string]any{
		"orderID":   orderID,
		"fromState": nil,
		"toState":   enums.Accepted,
		"timestamp": time.Now(),
	})
	s.publishEvent(event)

	return nil
}

func (s *OrderService) GetOrder(orderID string) (*dto.OrderCtx, error) {
	ctx, err := s.repo.GetOrder(orderID)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (s *OrderService) AcceptOrder(orderID string) error {
	return s.transitionOrder(orderID, acceptOrder)
}

func (s *OrderService) PrepareOrder(orderID string) error {
	return s.transitionOrder(orderID, prepareOrder)
}

func (s *OrderService) ReadyOrder(orderID string) error {
	return s.transitionOrder(orderID, readyOrder)
}

func (s *OrderService) PickOrder(orderID string) error {
	return s.transitionOrder(orderID, pickOrder)
}

func (s *OrderService) DeliverOrder(orderID string) error {
	return s.transitionOrder(orderID, deliverOrder)
}

func (s *OrderService) CancelOrder(orderID string) error {
	return s.transitionOrder(orderID, cancelOrder)
}

func (s *OrderService) transitionOrder(orderID string, action orderAction) error {
	fromState, toState := transitionStates(action)
	ctx, err := s.repo.GetOrder(orderID)
	if err != nil {
		return err
	}

	if ok, err := s.validateTransition(ctx.CurrentState, toState); !ok {
		return err
	}

	cmd := NewTransitionCommand(orderID, fromState, toState, s.auditRepo)

	if err := s.repo.TransitionTo(orderID, toState); err != nil {
		return err
	}
	cmd.Execute()

	fmt.Printf("Order transitioned successfully: %s", ctx.OrderID)
	event := utils.NewEvent(string(toState), map[string]any{
		"orderID":   orderID,
		"fromState": fromState,
		"toState":   toState,
		"timestamp": time.Now(),
	})
	s.publishEvent(event)

	return nil
}

func (s *OrderService) GetTransitionHistory(orderID string) []TransitionAudit {
	return s.auditRepo.History(orderID)
}

func (s *OrderService) validateTransition(fromState, toState enums.State) (bool, error) {
	validTransitions := map[enums.State][]enums.State{
		enums.Placed: {
			enums.Accepted,
			enums.Cancelled,
		},
		enums.Accepted: {
			enums.Preparing,
		},
		enums.Preparing: {
			enums.Ready,
		},
		enums.Ready: {
			enums.PickedUp,
		},
		enums.PickedUp: {
			enums.Delivered,
		},
	}

	if slices.Contains(validTransitions[fromState], toState) {
		return true, nil
	}

	return false, InvalidOrderTransition{
		FromState: fromState,
		ToState:   toState,
	}
}

func transitionStates(action orderAction) (enums.State, enums.State) {
	transitions := map[orderAction]struct {
		fromState enums.State
		toState   enums.State
	}{
		acceptOrder: {
			fromState: enums.Placed,
			toState:   enums.Accepted,
		},
		prepareOrder: {
			fromState: enums.Accepted,
			toState:   enums.Preparing,
		},
		readyOrder: {
			fromState: enums.Preparing,
			toState:   enums.Ready,
		},
		pickOrder: {
			fromState: enums.Ready,
			toState:   enums.PickedUp,
		},
		deliverOrder: {
			fromState: enums.PickedUp,
			toState:   enums.Delivered,
		},
		cancelOrder: {
			fromState: enums.Placed,
			toState:   enums.Cancelled,
		},
	}

	transition := transitions[action]
	return transition.fromState, transition.toState
}

func (s *OrderService) publishEvent(event interfaces.Event) {
	go func(event interfaces.Event) {
		err := s.bus.Publish(event)
		if err != nil {
			fmt.Errorf("%w", err)
		}
	}(event)
}
