package main

import (
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/delivery"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/eventbus"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

func main() {
	bus := eventbus.New()

	// services
	paymentSvc := payment.NewPaymentService()
	deliverySvc := delivery.NewDeliveryService()
	notificationSvc := notification.NewNotificationService()

	// observers
	paymentObs := payment.NewPaymentObserver(paymentSvc)
	deliveryObs := delivery.NewDeliveryObserver(deliverySvc)
	notifyObs := notification.NewNotificationObserver(notificationSvc)

	// subscriptions
	acceptOrder := interfaces.NewEvent(string(enums.Accepted), nil)
	cancelOrder := interfaces.NewEvent(string(enums.Cancelled), nil)
	readyEvent := interfaces.NewEvent(string(enums.Ready), nil)

	bus.Attach(acceptOrder, paymentObs)
	bus.Attach(cancelOrder, paymentObs)
	bus.Attach(readyEvent, deliveryObs)
	bus.Attach(acceptOrder, notifyObs)

	repo := order.NewOrderRepository()
	orderSvc := order.NewOrderService(repo, bus)

	orderSvc.PlaceOrder("ORD-001")
	orderSvc.AcceptOrder("ORD-001")
}
