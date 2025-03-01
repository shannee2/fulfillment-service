package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"fulfillment/dao"
	"fulfillment/kafka/producer"
	"fulfillment/model"
	"log"
	"math"
)

type OrderService struct {
	partnerService *DeliveryPartnerService
	assignmentDao  *dao.DeliveryAssignmentDAO
	partnerDao     *dao.DeliveryPartnerDAO
	producer       *producer.DeliveryAssignmentProducer
}

func NewOrderService(partnerService *DeliveryPartnerService, assignmentDao *dao.DeliveryAssignmentDAO, partnerDao *dao.DeliveryPartnerDAO, producer *producer.DeliveryAssignmentProducer) *OrderService {
	return &OrderService{
		partnerService: partnerService,
		assignmentDao:  assignmentDao,
		partnerDao:     partnerDao,
		producer:       producer,
	}
}

func (s *OrderService) ProcessOrder(orderJson []byte) {
	var order model.Order
	if err := json.Unmarshal(orderJson, &order); err != nil {
		log.Printf("Failed to parse order JSON: %v", err)
		return
	}

	fmt.Printf("Processing Order: %+v\n", order)

	partners, err := s.partnerService.GetPartnersByCity(order.Restaurant.City)
	if err != nil || len(partners) == 0 {
		log.Println("No available delivery partners in the city")
		return
	}

	nearestPartner, err := s.findNearestPartner(order.Restaurant.Latitude, order.Restaurant.Longitude, partners)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Assigned Delivery Partner: %s (ID: %d)\n", nearestPartner.Name, nearestPartner.ID)

	err = s.partnerDao.UpdateDeliveryPartnerStatus(nearestPartner.ID, "DELIVERING_ORDER")
	if err != nil {
		return
	}

	// Insert delivery assignment into DB
	assignment := model.DeliveryAssignment{
		DeliveryPartnerID: nearestPartner.ID,
		OrderID:           order.OrderID,
	}
	_, err = s.assignmentDao.InsertDeliveryAssignment(assignment)
	if err != nil {
		log.Printf("Failed to insert delivery assignment: %v", err)
		return
	}

	// Send Kafka message
	err = s.producer.ProduceAssignmentMessage(order.OrderID, nearestPartner.ID)
	if err != nil {
		log.Printf("Failed to produce Kafka message: %v", err)
	}
}

func (s *OrderService) findNearestPartner(orderLat, orderLon float64, partners []model.DeliveryPartner) (model.DeliveryPartner, error) {
	if len(partners) == 0 {
		return model.DeliveryPartner{}, errors.New("no delivery partners available")
	}

	var nearest model.DeliveryPartner
	minDistance := math.MaxFloat64

	for _, partner := range partners {
		lat, lon, _ := s.partnerService.GetLocation(partner.ID)
		distance := FindDistanceInMeters(orderLat, orderLon, lat, lon)

		if distance < minDistance {
			minDistance = distance
			nearest = partner
		}
	}

	if nearest.ID == 0 {
		return model.DeliveryPartner{}, errors.New("no nearby partner found")
	}

	fmt.Println("Distance to nearest partner:", minDistance)

	return nearest, nil
}
