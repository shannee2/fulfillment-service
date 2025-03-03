package delivery_assignment

import (
	"errors"
	"fulfillment/dao"
	"fulfillment/model"
)

type DeliveryAssignmentService struct {
	assignmentDAO dao.DeliveryAssignmentDAOInterface // Use Interface
}

func NewDeliveryAssignmentService(assignmentDAO dao.DeliveryAssignmentDAOInterface) *DeliveryAssignmentService {
	return &DeliveryAssignmentService{assignmentDAO: assignmentDAO}
}

func (s *DeliveryAssignmentService) AssignOrderToPartner(assignment model.DeliveryAssignment) (int, error) {
	// Validate input
	if assignment.DeliveryPartnerID == 0 || assignment.OrderID == 0 {
		return 0, errors.New("invalid assignment: OrderID and DeliveryPartnerID must be non-zero")
	}

	return s.assignmentDAO.InsertDeliveryAssignment(assignment)
}
