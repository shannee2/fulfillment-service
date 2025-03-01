package service

import (
	"fulfillment/dao"
	"fulfillment/model"
)

type DeliveryAssignmentService struct {
	assignmentDAO *dao.DeliveryAssignmentDAO
}

func NewDeliveryAssignmentService(assignmentDAO *dao.DeliveryAssignmentDAO) *DeliveryAssignmentService {
	return &DeliveryAssignmentService{assignmentDAO: assignmentDAO}
}

// Assigns an order to a delivery partner
func (s *DeliveryAssignmentService) AssignOrderToPartner(assignment model.DeliveryAssignment) (int, error) {
	return s.assignmentDAO.InsertDeliveryAssignment(assignment)
}
