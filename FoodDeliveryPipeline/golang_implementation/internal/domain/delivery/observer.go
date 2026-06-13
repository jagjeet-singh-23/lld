package delivery

import (
	"fmt"
	"sync"

	orderStates "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
	shared "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

type DeliveryObserver struct {
	service *DeliveryService
	mu      sync.RWMutex
}

func NewDeliveryObserver(service *DeliveryService) *DeliveryObserver {
	return &DeliveryObserver{
		service: service,
	}
}

func (o *DeliveryObserver) Update(event shared.Event) error {
	o.mu.RLock()
	orderID, err := validateMetadata(event.Metadata)
	if err != nil {
		o.mu.Unlock()
		return err
	}

	if event.EventType != string(orderStates.Ready) {
		o.mu.Unlock()
		return OrderNotReadyYet{
			OrderID: orderID,
		}
	}
	o.mu.Unlock()

	o.mu.Lock()
	defer o.mu.Unlock()
	partnerID, err := o.service.repo.Assign(orderID)
	if err != nil {
		return err
	}

	fmt.Printf("Delivery Parnter: %s assigned order: %s", partnerID, orderID)
	return nil
}

func validateMetadata(metadata map[string]any) (string, error) {
	orderID, ok := metadata["orderID"].(string)
	if !ok {
		return "", fmt.Errorf("invalid orderID")
	}

	return orderID, nil
}
