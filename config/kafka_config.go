package config

import "github.com/confluentinc/confluent-kafka-go/kafka"

func NewKafkaConsumer(broker string, groupID string) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}
	return kafka.NewConsumer(config)
}
