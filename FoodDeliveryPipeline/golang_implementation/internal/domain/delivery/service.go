package delivery

type DeliveryService struct {
	repo *DeliveryRepository
}

func NewDeliveryService() *DeliveryService {
	return &DeliveryService{}
}
