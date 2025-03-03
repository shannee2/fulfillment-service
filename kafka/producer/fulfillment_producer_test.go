package producer

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafkaProducer is a mock implementation of Kafka Producer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	args := m.Called(msg, deliveryChan)
	return args.Error(0)
}

func TestProduceAssignmentMessage_Success(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	topic := "test-topic"
	producer := &DeliveryAssignmentProducer{
		Producer: mockProducer,
		Topic:    topic,
	}

	orderID := int64(12345)
	deliveryPartnerID := 678

	expectedMessage := map[string]interface{}{
		"order_id":            orderID,
		"delivery_partner_id": deliveryPartnerID,
	}

	msgBytes, _ := json.Marshal(expectedMessage)

	mockProducer.On("Produce", mock.MatchedBy(func(msg *kafka.Message) bool {
		return string(msg.Value) == string(msgBytes) &&
			*msg.TopicPartition.Topic == topic
	}), mock.Anything).Return(nil)

	err := producer.ProduceAssignmentMessage(orderID, deliveryPartnerID)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestProduceAssignmentMessage_ProduceError(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	topic := "test-topic"
	producer := &DeliveryAssignmentProducer{
		Producer: mockProducer,
		Topic:    topic,
	}

	orderID := int64(12345)
	deliveryPartnerID := 678

	expectedMessage := map[string]interface{}{
		"order_id":            orderID,
		"delivery_partner_id": deliveryPartnerID,
	}

	msgBytes, _ := json.Marshal(expectedMessage)

	mockProducer.On("Produce", mock.MatchedBy(func(msg *kafka.Message) bool {
		return string(msg.Value) == string(msgBytes) &&
			*msg.TopicPartition.Topic == topic
	}), mock.Anything).Return(errors.New("failed to produce message"))

	err := producer.ProduceAssignmentMessage(orderID, deliveryPartnerID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to produce message")

	mockProducer.AssertExpectations(t)
}
