package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type DeliveryAssignmentProducerInterface interface {
	ProduceAssignmentMessage(orderID int64, deliveryPartnerID int) error
}

type KafkaProducerInterface interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

type DeliveryAssignmentProducer struct {
	Producer KafkaProducerInterface
	Topic    string
}

func NewDeliveryAssignmentProducer(broker, topic string) (*DeliveryAssignmentProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &DeliveryAssignmentProducer{Producer: p, Topic: topic}, nil
}

// ProduceAssignmentMessage produces a Kafka message
func (p *DeliveryAssignmentProducer) ProduceAssignmentMessage(orderID int64, deliveryPartnerID int) error {
	message := map[string]interface{}{
		"order_id":            orderID,
		"delivery_partner_id": deliveryPartnerID,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	err = p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          msgBytes,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	log.Printf("Produced assignment: %s\n", string(msgBytes))
	return nil
}
