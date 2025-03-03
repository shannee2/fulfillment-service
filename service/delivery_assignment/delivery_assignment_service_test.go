package delivery_assignment

import (
	"errors"
	"fulfillment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Mock DAO
type MockAssignmentDAO struct {
	mock.Mock
}

// Implement InsertDeliveryAssignment
func (m *MockAssignmentDAO) InsertDeliveryAssignment(assignment model.DeliveryAssignment) (int, error) {
	args := m.Called(assignment)
	return args.Int(0), args.Error(1)
}

func TestAssignOrderToPartner_Success(t *testing.T) {
	mockDAO := new(MockAssignmentDAO)
	service := NewDeliveryAssignmentService(mockDAO)

	assignment := model.DeliveryAssignment{
		DeliveryPartnerID: 123,
		OrderID:           456,
	}

	mockDAO.On("InsertDeliveryAssignment", assignment).Return(1, nil)

	id, err := service.AssignOrderToPartner(assignment)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	mockDAO.AssertExpectations(t)
}

func TestAssignOrderToPartner_InvalidAssignment(t *testing.T) {
	mockDAO := new(MockAssignmentDAO)
	service := NewDeliveryAssignmentService(mockDAO)

	invalidAssignment := model.DeliveryAssignment{
		DeliveryPartnerID: 0, // Invalid ID
		OrderID:           456,
	}

	id, err := service.AssignOrderToPartner(invalidAssignment)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "invalid assignment: OrderID and DeliveryPartnerID must be non-zero", err.Error())
}

func TestAssignOrderToPartner_DAOFailure(t *testing.T) {
	mockDAO := new(MockAssignmentDAO)
	service := NewDeliveryAssignmentService(mockDAO)

	assignment := model.DeliveryAssignment{
		DeliveryPartnerID: 123,
		OrderID:           456,
	}

	mockDAO.On("InsertDeliveryAssignment", assignment).Return(0, errors.New("DB error"))

	id, err := service.AssignOrderToPartner(assignment)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "DB error", err.Error())
	mockDAO.AssertExpectations(t)
}
