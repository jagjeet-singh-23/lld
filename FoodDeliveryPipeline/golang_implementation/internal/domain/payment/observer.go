package payment

import (
	"fmt"

	orderStates "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/enums"
	shared "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

type PaymentObserver struct {
	service *PaymentService
}

type PaymentMetadata struct {
	OrderID       string
	Amount        float32
	Payee         string
	Recipient     string
	ModeOfPayment enums.PaymentChannel
}

func NewPaymentObserver(service *PaymentService) *PaymentObserver {
	return &PaymentObserver{
		service: service,
	}
}

func (o *PaymentObserver) Update(event shared.Event) error {
	data, err := validateMetadata(event.Metadata)
	if err != nil {
		return err
	}

	if event.EventType == string(orderStates.Accepted) {
		err := o.service.CapturePayment(
			data.OrderID,
			data.Amount,
			data.Payee,
			data.Recipient,
			data.ModeOfPayment,
		)
		if err != nil {
			return err
		}

		return nil
	}

	if event.EventType == string(orderStates.Cancelled) {
		err := o.service.Refund(data.OrderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateMetadata(metadata map[string]any) (*PaymentMetadata, error) {
	orderID, ok := metadata["orderID"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid orderID")
	}

	amount, ok := metadata["amount"].(float32)
	if !ok {
		return nil, fmt.Errorf("invalid amount")
	}

	payee, ok := metadata["payee"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid payee")
	}

	recipient, ok := metadata["recepient"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid recipient")
	}

	modeOfPayment, ok := metadata["modeOfPayment"].(enums.PaymentChannel)
	if !ok {
		return nil, fmt.Errorf("invalid modeOfPayment")
	}

	return &PaymentMetadata{
		OrderID:       orderID,
		Amount:        amount,
		Payee:         payee,
		Recipient:     recipient,
		ModeOfPayment: modeOfPayment,
	}, nil
}
