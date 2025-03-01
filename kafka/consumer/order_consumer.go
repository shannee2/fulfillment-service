package consumer

import (
	"fmt"
	"fulfillment/config"
	"fulfillment/service"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderConsumer struct {
	consumer     *kafka.Consumer
	orderService *service.OrderService
}

func NewOrderConsumer(broker string, groupID string, orderService *service.OrderService) (*OrderConsumer, error) {
	consumer, err := config.NewKafkaConsumer(broker, groupID)
	if err != nil {
		return nil, err
	}

	return &OrderConsumer{
		consumer:     consumer,
		orderService: orderService,
	}, nil
}

func (c *OrderConsumer) Start(topic string) {
	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	fmt.Println("Kafka Order Consumer is running...")

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", msg.Value)
			c.orderService.ProcessOrder(msg.Value)
		} else {
			log.Printf("Consumer error: %v", err)
		}
	}
}
