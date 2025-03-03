package delivery_partner

import (
	//"fulfillment/dao"
	"fulfillment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Mock DAO
// Mock DAO implementing the interface
type MockPartnerDAO struct {
	mock.Mock
}

// Implementing methods from DeliveryPartnerDAOInterface
func (m *MockPartnerDAO) InsertDeliveryPartner(partner model.DeliveryPartner) error {
	args := m.Called(partner)
	return args.Error(0)
}

func (m *MockPartnerDAO) UpdateDeliveryPartnerStatus(partnerID int, status string) error {
	args := m.Called(partnerID, status)
	return args.Error(0)
}

func (m *MockPartnerDAO) GetAvailableDeliveryPartnersByCity(city string) ([]model.DeliveryPartner, error) {
	args := m.Called(city)
	return args.Get(0).([]model.DeliveryPartner), args.Error(1)
}

func TestAddDeliveryPartner(t *testing.T) {
	mockDAO := new(MockPartnerDAO)
	service := NewDeliveryPartnerService(mockDAO)

	partner := model.DeliveryPartner{ID: 1, Name: "John Doe", City: "New York"}

	mockDAO.On("InsertDeliveryPartner", partner).Return(nil)

	err := service.AddDeliveryPartner(partner)

	assert.NoError(t, err)
	mockDAO.AssertExpectations(t)
}

func TestUpdatePartnerStatus(t *testing.T) {
	mockDAO := new(MockPartnerDAO)
	service := NewDeliveryPartnerService(mockDAO)

	partnerID := 1
	status := model.DeliveryStatus("AVAILABLE")

	mockDAO.On("UpdateDeliveryPartnerStatus", partnerID, string(status)).Return(nil)

	err := service.UpdatePartnerStatus(partnerID, status)

	assert.NoError(t, err)
	mockDAO.AssertExpectations(t)
}

func TestGetPartnersByCity_Success(t *testing.T) {
	mockDAO := new(MockPartnerDAO)
	service := NewDeliveryPartnerService(mockDAO)

	city := "New York"
	expectedPartners := []model.DeliveryPartner{
		{ID: 1, Name: "John Doe", City: city},
		{ID: 2, Name: "Jane Doe", City: city},
	}

	mockDAO.On("GetAvailableDeliveryPartnersByCity", city).Return(expectedPartners, nil)

	partners, err := service.GetPartnersByCity(city)

	assert.NoError(t, err)
	assert.Equal(t, expectedPartners, partners)
	mockDAO.AssertExpectations(t)
}

func TestGetLocation(t *testing.T) {
	service := NewDeliveryPartnerService(nil)

	lat, lon, err := service.GetLocation(1)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, lat, 20.0)
	assert.LessOrEqual(t, lat, 30.0)
	assert.GreaterOrEqual(t, lon, 75.0)
	assert.LessOrEqual(t, lon, 85.0)
}
