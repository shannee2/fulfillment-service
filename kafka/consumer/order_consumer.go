package consumer

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaConsumerInterface abstracts the Kafka consumer methods
type KafkaConsumerInterface interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
}

// OrderServiceInterface abstracts the order fulfillment service
type OrderServiceInterface interface {
	ProcessOrder(orderJson []byte)
}

type OrderConsumer struct {
	Consumer     KafkaConsumerInterface
	OrderService OrderServiceInterface
}

func NewOrderConsumer(consumer KafkaConsumerInterface, orderService OrderServiceInterface) *OrderConsumer {
	return &OrderConsumer{
		Consumer:     consumer,
		OrderService: orderService,
	}
}

func (c *OrderConsumer) Start(topic string) {
	err := c.Consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	fmt.Println("Kafka Order Consumer is running...")

	for {
		msg, err := c.Consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", msg.Value)
			c.OrderService.ProcessOrder(msg.Value)
		} else {
			log.Printf("Consumer error: %v", err)
		}
	}
}
