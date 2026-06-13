package dto

import (
	"time"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/enums"
)

type PaymentContext struct {
	TxnID         string
	MerchantTxnID string
	OrderID       string
	Amount        float32
	Payee         string
	Recepient     string
	Mode          enums.PaymentChannel
	Timestamp     time.Time

	State enums.PaymentState
}

func NewPaymentContext(txnID string, merchantTxnID string, orderID string, amount float32, payee string, recepient string, modeOfPayment enums.PaymentChannel) *PaymentContext {
	return &PaymentContext{
		TxnID:         txnID,
		MerchantTxnID: merchantTxnID,
		OrderID:       orderID,
		Amount:        amount,
		Payee:         payee,
		Recepient:     recepient,
		Mode:          modeOfPayment,
		State:         enums.Processing,
		Timestamp:     time.Now(),
	}
}
