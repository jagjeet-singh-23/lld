package payment

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/dto"
)

type PaymentError interface {
	Error() string
}

type InsufficientBalance struct {
	ctx dto.PaymentContext
}

func (e InsufficientBalance) Error() string {
	return fmt.Sprintf("Unable to process payment due to insufficent balance, %s, requested amount: %0.2f", e.ctx.TxnID, e.ctx.Amount)
}
