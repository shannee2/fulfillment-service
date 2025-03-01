package main

import (
	"fulfillment/dao"
	"fulfillment/kafka/consumer"
	"fulfillment/kafka/producer"
	"fulfillment/service"
	"log"
)

func main() {
	// Initialize database connection
	db, err := dao.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Ensure DB connection closes on exit

	// Kafka configuration
	broker := "localhost:9092"
	topic := "orders"
	groupID := "order-consumer-group"

	// Initialize DAOs
	deliveryPartnerDAO := dao.NewDeliveryPartnerDAO(db)
	deliveryAssignmentDAO := dao.NewDeliveryAssignmentDAO(db)

	// Initialize Services
	deliveryPartnerService := service.NewDeliveryPartnerService(deliveryPartnerDAO)
	producer, _ := producer.NewDeliveryAssignmentProducer(broker, "delivery_assignment")
	orderService := service.NewOrderService(deliveryPartnerService, deliveryAssignmentDAO, deliveryPartnerDAO, producer)

	// Initialize Kafka consumer
	consumer, err := consumer.NewOrderConsumer(broker, groupID, orderService)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	// Start consuming messages
	consumer.Start(topic)
}
