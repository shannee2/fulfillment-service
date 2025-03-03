package delivery_partner

import (
	"fulfillment/dao"
	"fulfillment/model"
	"math/rand"
)

type DeliveryPartnerService interface {
	AddDeliveryPartner(partner model.DeliveryPartner) error
	UpdatePartnerStatus(partnerID int, status model.DeliveryStatus) error
	GetPartnersByCity(city string) ([]model.DeliveryPartner, error)
	GetLocation(partnerID int) (float64, float64, error)
}

type DeliveryPartnerServiceImpl struct {
	dao dao.DeliveryPartnerDAOInterface // Use the interface, not concrete struct
}

func NewDeliveryPartnerService(dao dao.DeliveryPartnerDAOInterface) *DeliveryPartnerServiceImpl {
	return &DeliveryPartnerServiceImpl{dao: dao}
}

func (s *DeliveryPartnerServiceImpl) AddDeliveryPartner(partner model.DeliveryPartner) error {
	return s.dao.InsertDeliveryPartner(partner)
}

func (s *DeliveryPartnerServiceImpl) UpdatePartnerStatus(partnerID int, status model.DeliveryStatus) error {
	return s.dao.UpdateDeliveryPartnerStatus(partnerID, string(status))
}

func (s *DeliveryPartnerServiceImpl) GetPartnersByCity(city string) ([]model.DeliveryPartner, error) {
	return s.dao.GetAvailableDeliveryPartnersByCity(city)
}

func (s *DeliveryPartnerServiceImpl) GetLocation(partnerID int) (float64, float64, error) {
	lat := 20.0 + rand.Float64()*10
	lon := 75.0 + rand.Float64()*10
	return lat, lon, nil
}
