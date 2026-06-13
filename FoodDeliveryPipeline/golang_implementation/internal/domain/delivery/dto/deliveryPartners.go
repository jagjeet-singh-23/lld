package dto

type DeliveryPartner struct {
	PartnerID       string
	CurrentLocation string
}

func NewDeliveryParnter(partnerID string, currentLocation string) *DeliveryPartner {
	return &DeliveryPartner{
		PartnerID:       partnerID,
		CurrentLocation: currentLocation,
	}
}
