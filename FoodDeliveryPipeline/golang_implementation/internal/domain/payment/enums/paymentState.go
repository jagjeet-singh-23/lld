package enums

type PaymentState string

const (
	Processing PaymentState = "processing"
	Successful PaymentState = "success"
	Failed     PaymentState = "failed"
)
