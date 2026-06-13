package payment

import (
	"fmt"
	"sync"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/dto"
)

type PaymentRepository struct {
	payments         map[string]*dto.PaymentContext
	paymentsCaptured []*dto.PaymentContext
	paymentsRefunded []*dto.PaymentContext
	mu               sync.RWMutex
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		payments:         make(map[string]*dto.PaymentContext),
		paymentsCaptured: make([]*dto.PaymentContext, 0),
		paymentsRefunded: make([]*dto.PaymentContext, 0),
	}
}

func (r *PaymentRepository) CapturePayment(ctx *dto.PaymentContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.addTxn(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepository) GetPayment(orderID string) (*dto.PaymentContext, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.hasTxn(orderID)
}

func (r *PaymentRepository) Refund(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx, err := r.hasTxn(orderID)
	if err != nil {
		return err
	}

	if err := r.addRefund(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) hasTxn(orderID string) (*dto.PaymentContext, error) {
	if _, ok := r.payments[orderID]; !ok {
		return nil, fmt.Errorf("no such txn found: %s", orderID)
	}

	return r.payments[orderID], nil
}

func (r *PaymentRepository) addTxn(ctx *dto.PaymentContext) error {
	if _, err := r.hasTxn(ctx.OrderID); err == nil {
		return fmt.Errorf("duplicate txn: %s", ctx.OrderID)
	}

	r.paymentsCaptured = append(r.paymentsCaptured, ctx)
	r.payments[ctx.OrderID] = ctx

	return nil
}

func (r *PaymentRepository) addRefund(ctx *dto.PaymentContext) error {
	for idx, payment := range r.paymentsCaptured {
		if payment.OrderID == ctx.OrderID {
			r.paymentsCaptured = append(r.paymentsCaptured[:idx], r.paymentsCaptured[idx+1:]...)
			break
		}
	}

	r.paymentsRefunded = append(r.paymentsRefunded, ctx)
	return nil
}
