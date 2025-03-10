package consumer

import (
	"fmt"
	"fulfillment/config"
	"fulfillment/service/order_fulfillment"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderConsumer struct {
	Consumer     *kafka.Consumer
	OrderService *order_fulfillment.OrderService
}

func NewOrderConsumer(broker string, groupID string, orderService *order_fulfillment.OrderService) (*OrderConsumer, error) {
	consumer, err := config.NewKafkaConsumer(broker, groupID)
	if err != nil {
		return nil, err
	}

	return &OrderConsumer{
		Consumer:     consumer,
		OrderService: orderService,
	}, nil
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
