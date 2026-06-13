package enums

type PaymentChannel string

const (
	Upi        PaymentChannel = "upi"
	CreditCard PaymentChannel = "creditCard"
	COD        PaymentChannel = "cashOnDelivery"
)
