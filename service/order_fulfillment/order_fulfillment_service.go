package order_fulfillment

import (
	"encoding/json"
	"errors"
	"fmt"
	"fulfillment/dao"
	"fulfillment/kafka/producer"
	"fulfillment/model"
	"fulfillment/service/delivery_partner"
	"fulfillment/service/location"
	"log"
	"math"
)

type OrderService struct {
	partnerService delivery_partner.DeliveryPartnerService
	assignmentDao  dao.DeliveryAssignmentDAOInterface
	partnerDao     dao.DeliveryPartnerDAOInterface
	producer       producer.DeliveryAssignmentProducerInterface
}

func NewOrderService(
	partnerService delivery_partner.DeliveryPartnerService,
	assignmentDao dao.DeliveryAssignmentDAOInterface,
	partnerDao dao.DeliveryPartnerDAOInterface,
	producer producer.DeliveryAssignmentProducerInterface,
) *OrderService {
	return &OrderService{
		partnerService: partnerService,
		assignmentDao:  assignmentDao,
		partnerDao:     partnerDao,
		producer:       producer,
	}
}

func (s *OrderService) ProcessOrder(orderJson []byte) {
	// Parse the order JSON
	order, err := s.parseOrder(orderJson)
	if err != nil {
		log.Printf("Failed to parse order JSON: %v", err)
		return
	}

	// Assign a delivery partner
	nearestPartner, err := s.assignDeliveryPartner(order)
	if err != nil {
		log.Println(err)
		return
	}

	// Update partner status
	if err := s.partnerDao.UpdateDeliveryPartnerStatus(nearestPartner.ID, "DELIVERING_ORDER"); err != nil {
		log.Printf("Failed to update delivery partner status: %v", err)
		return
	}

	// Insert delivery assignment
	assignment := model.DeliveryAssignment{
		DeliveryPartnerID: nearestPartner.ID,
		OrderID:           order.OrderID,
	}
	if _, err := s.assignmentDao.InsertDeliveryAssignment(assignment); err != nil {
		log.Printf("Failed to insert delivery assignment: %v", err)
		return
	}

	// Produce Kafka message
	if err := s.producer.ProduceAssignmentMessage(order.OrderID, nearestPartner.ID); err != nil {
		log.Printf("Failed to produce Kafka message: %v", err)
	}
}

// Helper function to parse order JSON
func (s *OrderService) parseOrder(orderJson []byte) (model.Order, error) {
	var order model.Order
	if err := json.Unmarshal(orderJson, &order); err != nil {
		return model.Order{}, err
	}
	fmt.Printf("Processing Order: %+v\n", order)
	return order, nil
}

// Helper function to assign a delivery partner
func (s *OrderService) assignDeliveryPartner(order model.Order) (model.DeliveryPartner, error) {
	partners, err := s.partnerService.GetPartnersByCity(order.Restaurant.City)
	if err != nil || len(partners) == 0 {
		return model.DeliveryPartner{}, errors.New("no available delivery partners in the city")
	}

	nearestPartner, err := s.findNearestPartner(order.Restaurant.Latitude, order.Restaurant.Longitude, partners)
	if err != nil {
		return model.DeliveryPartner{}, err
	}

	fmt.Printf("Assigned Delivery Partner: %s (ID: %d)\n", nearestPartner.Name, nearestPartner.ID)
	return nearestPartner, nil
}

func (s *OrderService) findNearestPartner(orderLat, orderLon float64, partners []model.DeliveryPartner) (model.DeliveryPartner, error) {
	if len(partners) == 0 {
		return model.DeliveryPartner{}, errors.New("no delivery partners available")
	}

	var nearest model.DeliveryPartner
	minDistance := math.MaxFloat64

	for _, partner := range partners {
		lat, lon, _ := s.partnerService.GetLocation(partner.ID)
		distance := location.FindDistanceInMeters(orderLat, orderLon, lat, lon)

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
