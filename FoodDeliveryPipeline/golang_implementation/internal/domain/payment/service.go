package payment

import (
	"fmt"
	"sync"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/dto"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/payment/gateway"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/utils"
)

type PaymentService struct {
	repo     *PaymentRepository
	registry *gateway.PaymentGatewayRegistry
	mu       sync.RWMutex
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		repo:     NewPaymentRepository(),
		registry: gateway.NewPaymentGatewayRegistry(),
	}
}

func (s *PaymentService) CapturePayment(orderID string, amount float32, payee, recepient string, modeOfPayment enums.PaymentChannel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("Payment capture initiated: %s, payee: %s, recepient: %s\n", orderID, payee, recepient)

	// 0. Idempotency check: has this order already been paid or is it currently processing?
	if existing, err := s.repo.GetPayment(orderID); err == nil && existing != nil {
		if existing.State == enums.Successful {
			fmt.Printf("Payment for order %s already captured successfully (Idempotent call). TxnID: %s, MerchantTxnID: %s\n", orderID, existing.TxnID, existing.MerchantTxnID)
			return nil
		}
		if existing.State == enums.Processing {
			return fmt.Errorf("payment for order %s is already processing", orderID)
		}
	}

	// 1. Generate our internal transaction ID (our TxnID)
	txnID, err := utils.GenerateTxnID()
	if err != nil {
		return err
	}

	// 2. Prepare the initial payment context
	ctx := s.preparePaymentContext(txnID, orderID, amount, payee, recepient, modeOfPayment)

	// 3. Resolve the dynamic payment gateway using the ModeOfPayment
	gw, err := s.registry.GetGateway(modeOfPayment)
	if err != nil {
		return err
	}

	// 4. Process the payment via external gateway and receive the merchant transaction ID
	merchantTxnID, err := gw.Process(*ctx)
	if err != nil {
		ctx.State = enums.Failed
		_ = s.repo.CapturePayment(ctx) // Persist failed attempt
		return err
	}

	// 5. Update the context with merchant's transaction ID and successful state
	ctx.MerchantTxnID = merchantTxnID
	ctx.State = enums.Successful

	// 6. Capture/Persist the completed payment in our repository
	err = s.repo.CapturePayment(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Payment successful, TxnID: %s, MerchantTxnID: %s\n", txnID, merchantTxnID)
	return nil
}

func (s *PaymentService) preparePaymentContext(
	txnID string,
	orderID string,
	amount float32,
	payee string,
	recepient string,
	modeOfPayment enums.PaymentChannel,
) *dto.PaymentContext {
	return dto.NewPaymentContext(txnID, "", orderID, amount, payee, recepient, modeOfPayment)
}

func (s *PaymentService) Refund(orderID string) error {
	fmt.Printf("Refund initiated for order: %s\n", orderID)
	return s.repo.Refund(orderID)
}
