package consumer_test

//
//import (
//	"fulfillment/kafka/consumer"
//	"testing"
//	"time"
//
//	"github.com/confluentinc/confluent-kafka-go/kafka"
//	"github.com/stretchr/testify/mock"
//)
//
//// Mock Kafka Consumer
//type MockKafkaConsumer struct {
//	mock.Mock
//}
//
//func (m *MockKafkaConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
//	args := m.Called(topics, rebalanceCb)
//	return args.Error(0)
//}
//
//func (m *MockKafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
//	args := m.Called(timeout)
//	if msg, ok := args.Get(0).(*kafka.Message); ok {
//		return msg, args.Error(1)
//	}
//	return nil, args.Error(1)
//}
//
//// Mock OrderService
//type MockOrderService struct {
//	mock.Mock
//}
//
//func (m *MockOrderService) ProcessOrder(data []byte) {
//	m.Called(data)
//}
//
//func TestOrderConsumer_SuccessfulMessageProcessing(t *testing.T) {
//	mockConsumer := new(MockKafkaConsumer)
//	mockOrderService := new(MockOrderService)
//
//	orderConsumer := &consumer.OrderConsumer{
//		Consumer:     mockConsumer,
//		OrderService: mockOrderService,
//	}
//
//	topic := "test-topic"
//	message := []byte(`{"order_id": 12345, "status": "NEW"}`)
//
//	mockConsumer.On("SubscribeTopics", []string{topic}, mock.Anything).Return(nil)
//	mockConsumer.On("ReadMessage", mock.Anything).Return(&kafka.Message{Value: message}, nil).Maybe()
//	mockOrderService.On("ProcessOrder", message).Return(nil)
//
//	done := make(chan bool)
//
//	go func() {
//		orderConsumer.Start(topic)
//		done <- true
//	}()
//
//	// Wait for the consumer to process one message
//	<-done
//
//	// Validate expectations
//	mockConsumer.AssertExpectations(t)
//	mockOrderService.AssertExpectations(t)
//}
