package gateway

import (
	"fmt"
	"sync"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/dto"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/utils"
)

type UpiGateway struct {
	mu sync.RWMutex
}

func NewUpiGateway() *UpiGateway {
	return &UpiGateway{}
}

func (g *UpiGateway) Process(ctx dto.PaymentContext) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// mimicing a successful payment by generating a completely independent transaction ID
	merchantTxnID, err := utils.GenerateGatewayTxnID()
	if err != nil {
		return "", err
	}

	fmt.Printf("Payment successful, amount: %0.2f, payee: %s, recepient: %s, merchantTxnID: %s\n", ctx.Amount, ctx.Payee, ctx.Recepient, merchantTxnID)
	return merchantTxnID, nil
}
