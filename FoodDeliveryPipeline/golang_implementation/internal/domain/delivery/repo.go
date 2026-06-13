package delivery

import (
	"fmt"
	"sync"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/delivery/dto"
)

type DeliveryRepository struct {
	availableDeliveryPartners []dto.DeliveryPartner
	assignedDeliveryPartners  map[string]string
	mu                        sync.RWMutex
}

func NewDeliveryRepository() *DeliveryRepository {
	return &DeliveryRepository{
		availableDeliveryPartners: make([]dto.DeliveryPartner, 0),
		assignedDeliveryPartners:  make(map[string]string),
	}
}

func (r *DeliveryRepository) Assign(orderID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.canAssignDeliveryPartner(orderID)
	if err != nil {
		return "", err
	}

	partnerID, err := r.assignDeliveryPartner(orderID)
	if err != nil {
		return "", err
	}
	return partnerID, nil
}

func (r *DeliveryRepository) CompleteDelivery(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.freeDeliveryPartner(orderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeliveryRepository) canAssignDeliveryPartner(orderID string) error {
	if len(r.availableDeliveryPartners) == 0 {
		return NoDeliveryPartnerFound{
			OrderID: orderID,
		}
	}
	return nil
}

func (r *DeliveryRepository) assignDeliveryPartner(orderID string) (string, error) {
	deliveryPartner := r.availableDeliveryPartners[0]
	r.availableDeliveryPartners = r.availableDeliveryPartners[1:]
	r.assignedDeliveryPartners[orderID] = deliveryPartner.PartnerID
	return deliveryPartner.PartnerID, nil
}

func (r *DeliveryRepository) freeDeliveryPartner(orderID string) error {
	parnterID, ok := r.assignedDeliveryPartners[orderID]
	if !ok {
		return fmt.Errorf("no order found: %s", orderID)
	}

	r.availableDeliveryPartners = append(r.availableDeliveryPartners, *dto.NewDeliveryParnter(parnterID, ""))
	return nil
}
