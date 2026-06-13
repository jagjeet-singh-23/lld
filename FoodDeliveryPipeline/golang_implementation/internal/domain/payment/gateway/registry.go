package gateway

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/interfaces"
)

type PaymentGatewayRegistry struct {
	registry map[enums.PaymentChannel]interfaces.PaymentGateway
}

func NewPaymentGatewayRegistry() *PaymentGatewayRegistry {
	registry := map[enums.PaymentChannel]interfaces.PaymentGateway{
		enums.Upi: NewUpiGateway(),
	}
	return &PaymentGatewayRegistry{
		registry: registry,
	}
}

func (r *PaymentGatewayRegistry) GetGateway(channel enums.PaymentChannel) (interfaces.PaymentGateway, error) {
	gw, exists := r.registry[channel]
	if !exists {
		return nil, fmt.Errorf("no gateway registered for payment channel: %s", channel)
	}
	return gw, nil
}
