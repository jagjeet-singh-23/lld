package interfaces

import "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/dto"

type PaymentGateway interface {
	Process(ctx dto.PaymentContext) (string, error)
}
