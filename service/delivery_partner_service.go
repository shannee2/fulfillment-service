package service

import (
	"fulfillment/dao"
	"fulfillment/model"
	"math/rand"
)

type DeliveryPartnerService struct {
	dao *dao.DeliveryPartnerDAO
}

func NewDeliveryPartnerService(dao *dao.DeliveryPartnerDAO) *DeliveryPartnerService {
	return &DeliveryPartnerService{dao: dao}
}

// Insert a new delivery partner
func (s *DeliveryPartnerService) AddDeliveryPartner(partner model.DeliveryPartner) error {
	return s.dao.InsertDeliveryPartner(partner)
}

// Update delivery partner status
func (s *DeliveryPartnerService) UpdatePartnerStatus(partnerID int, status model.DeliveryStatus) error {
	return s.dao.UpdateDeliveryPartnerStatus(partnerID, string(status))
}

// Fetch all delivery partners in a specific city
func (s *DeliveryPartnerService) GetPartnersByCity(city string) ([]model.DeliveryPartner, error) {
	return s.dao.GetAvailableDeliveryPartnersByCity(city)
}

func (s *DeliveryPartnerService) GetLocation(partnerID int) (float64, float64, error) {
	//partners, err := s.dao.GetAvailableDeliveryPartnersByCity("any") // Assume any city for demo
	//if err != nil || len(partners) == 0 {
	//	return 0, 0, errors.New("no available partners")
	//}

	// Mocking random GPS coordinates
	lat := 20.0 + rand.Float64()*10
	lon := 75.0 + rand.Float64()*10
	return lat, lon, nil
}
