package order_fulfillment

import (
	"encoding/json"
	"errors"
	"fulfillment/model"
	"testing"

	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mocks ----

type MockDeliveryPartnerService struct {
	mock.Mock
}

func (m *MockDeliveryPartnerService) GetPartnersByCity(city string) ([]model.DeliveryPartner, error) {
	args := m.Called(city)
	if args.Get(0) != nil {
		return args.Get(0).([]model.DeliveryPartner), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeliveryPartnerService) GetLocation(partnerID int) (float64, float64, error) {
	args := m.Called(partnerID)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockDeliveryPartnerService) AddDeliveryPartner(partner model.DeliveryPartner) error {
	args := m.Called(partner)
	return args.Error(0)
}

func (m *MockDeliveryPartnerService) UpdatePartnerStatus(partnerID int, status model.DeliveryStatus) error {
	args := m.Called(partnerID, status)
	return args.Error(0)
}

type MockAssignmentDAO struct {
	mock.Mock
}

func (m *MockAssignmentDAO) InsertDeliveryAssignment(assignment model.DeliveryAssignment) (int, error) {
	args := m.Called(assignment)
	return args.Int(0), args.Error(1)
}

type MockPartnerDAO struct {
	mock.Mock
}

// Existing method
func (m *MockPartnerDAO) UpdateDeliveryPartnerStatus(partnerID int, status string) error {
	args := m.Called(partnerID, status)
	return args.Error(0)
}

// Missing method 1: InsertDeliveryPartner
func (m *MockPartnerDAO) InsertDeliveryPartner(partner model.DeliveryPartner) error {
	args := m.Called(partner)
	return args.Error(0)
}

// Missing method 2: GetAvailableDeliveryPartnersByCity
func (m *MockPartnerDAO) GetAvailableDeliveryPartnersByCity(city string) ([]model.DeliveryPartner, error) {
	args := m.Called(city)
	return args.Get(0).([]model.DeliveryPartner), args.Error(1)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) ProduceAssignmentMessage(orderID int64, deliveryPartnerID int) error {
	args := m.Called(orderID, deliveryPartnerID)
	return args.Error(0)
}

// ---- Test Cases ----

func TestProcessOrder_Success(t *testing.T) {
	mockPartnerService := new(MockDeliveryPartnerService)
	mockAssignmentDAO := new(MockAssignmentDAO)
	mockPartnerDAO := new(MockPartnerDAO)
	mockKafkaProducer := new(MockKafkaProducer)

	service := NewOrderService(mockPartnerService, mockAssignmentDAO, mockPartnerDAO, mockKafkaProducer)

	order := model.Order{
		OrderID: 1,
		Restaurant: model.Restaurant{
			City:      "New York",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
	}

	orderJSON, _ := json.Marshal(order)

	partners := []model.DeliveryPartner{
		{ID: 101, Name: "John Doe"},
	}

	// Mock responses
	mockPartnerService.On("GetPartnersByCity", "New York").Return(partners, nil)
	mockPartnerService.On("GetLocation", 101).Return(40.7130, -74.0070, nil)
	mockPartnerDAO.On("UpdateDeliveryPartnerStatus", 101, "DELIVERING_ORDER").Return(nil)
	mockAssignmentDAO.On("InsertDeliveryAssignment", mock.Anything).Return(1, nil)
	mockKafkaProducer.On("ProduceAssignmentMessage", int64(1), 101).Return(nil)

	// Execute
	service.ProcessOrder(orderJSON)

	// Assertions
	mockPartnerService.AssertExpectations(t)
	mockPartnerDAO.AssertExpectations(t)
	mockAssignmentDAO.AssertExpectations(t)
	mockKafkaProducer.AssertExpectations(t)
}

func TestProcessOrder_NoPartnersAvailable(t *testing.T) {
	mockPartnerService := new(MockDeliveryPartnerService)
	mockAssignmentDAO := new(MockAssignmentDAO)
	mockPartnerDAO := new(MockPartnerDAO)
	mockKafkaProducer := new(MockKafkaProducer)

	service := NewOrderService(mockPartnerService, mockAssignmentDAO, mockPartnerDAO, mockKafkaProducer)

	order := model.Order{
		OrderID: 2,
		Restaurant: model.Restaurant{
			City:      "Los Angeles",
			Latitude:  34.0522,
			Longitude: -118.2437,
		},
	}

	orderJSON, _ := json.Marshal(order)

	// Mock responses
	mockPartnerService.On("GetPartnersByCity", "Los Angeles").Return([]model.DeliveryPartner{}, nil)

	// Execute
	service.ProcessOrder(orderJSON)

	// Assertions
	mockPartnerService.AssertExpectations(t)
	mockAssignmentDAO.AssertNotCalled(t, "InsertDeliveryAssignment", mock.Anything)
	mockKafkaProducer.AssertNotCalled(t, "ProduceAssignmentMessage", mock.Anything)
}

func TestProcessOrder_DBInsertFailure(t *testing.T) {
	mockPartnerService := new(MockDeliveryPartnerService)
	mockAssignmentDAO := new(MockAssignmentDAO)
	mockPartnerDAO := new(MockPartnerDAO)
	mockKafkaProducer := new(MockKafkaProducer)

	service := NewOrderService(mockPartnerService, mockAssignmentDAO, mockPartnerDAO, mockKafkaProducer)

	order := model.Order{
		OrderID: 3,
		Restaurant: model.Restaurant{
			City:      "Chicago",
			Latitude:  41.8781,
			Longitude: -87.6298,
		},
	}

	orderJSON, _ := json.Marshal(order)

	partners := []model.DeliveryPartner{
		{ID: 102, Name: "Alice Smith"},
	}

	// Mock responses
	mockPartnerService.On("GetPartnersByCity", "Chicago").Return(partners, nil)
	mockPartnerService.On("GetLocation", 102).Return(41.8785, -87.6300, nil)
	mockPartnerDAO.On("UpdateDeliveryPartnerStatus", 102, "DELIVERING_ORDER").Return(nil)
	mockAssignmentDAO.On("InsertDeliveryAssignment", mock.Anything).Return(0, errors.New("DB error"))

	// Execute
	service.ProcessOrder(orderJSON)

	// Assertions
	mockPartnerService.AssertExpectations(t)
	mockPartnerDAO.AssertExpectations(t)
	mockAssignmentDAO.AssertExpectations(t)
	mockKafkaProducer.AssertNotCalled(t, "ProduceAssignmentMessage", mock.Anything)
}

func TestProcessOrder_KafkaFailure(t *testing.T) {
	mockPartnerService := new(MockDeliveryPartnerService)
	mockAssignmentDAO := new(MockAssignmentDAO)
	mockPartnerDAO := new(MockPartnerDAO)
	mockKafkaProducer := new(MockKafkaProducer)

	service := NewOrderService(mockPartnerService, mockAssignmentDAO, mockPartnerDAO, mockKafkaProducer)

	order := model.Order{
		OrderID: 4,
		Restaurant: model.Restaurant{
			City:      "Houston",
			Latitude:  29.7604,
			Longitude: -95.3698,
		},
	}

	orderJSON, _ := json.Marshal(order)

	partners := []model.DeliveryPartner{
		{ID: 103, Name: "Bob Johnson"},
	}

	// Mock responses
	mockPartnerService.On("GetPartnersByCity", "Houston").Return(partners, nil)
	mockPartnerService.On("GetLocation", 103).Return(29.7610, -95.3700, nil)
	mockPartnerDAO.On("UpdateDeliveryPartnerStatus", 103, "DELIVERING_ORDER").Return(nil)
	mockAssignmentDAO.On("InsertDeliveryAssignment", mock.Anything).Return(1, nil)
	mockKafkaProducer.On("ProduceAssignmentMessage", int64(4), 103).Return(errors.New("Kafka error"))

	// Execute
	service.ProcessOrder(orderJSON)

	// Assertions
	mockPartnerService.AssertExpectations(t)
	mockPartnerDAO.AssertExpectations(t)
	mockAssignmentDAO.AssertExpectations(t)
	mockKafkaProducer.AssertExpectations(t)
}
