package main

import (
	"fulfillment/dao"
	"fulfillment/kafka/consumer"
	"fulfillment/kafka/producer"
	"fulfillment/service/delivery_partner"
	"fulfillment/service/order_fulfillment"
	"log"
)

func main() {
	db, err := dao.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Kafka configuration
	broker := "localhost:9092"
	topic := "orders"
	groupID := "order-orderConsumer-group"

	// Initialize DAOs
	deliveryPartnerDAO := dao.NewDeliveryPartnerDAO(db)
	deliveryAssignmentDAO := dao.NewDeliveryAssignmentDAO(db)

	// Initialize Services
	deliveryPartnerService := delivery_partner.NewDeliveryPartnerService(deliveryPartnerDAO)
	deliveryAssignmentProducer, _ := producer.NewDeliveryAssignmentProducer(broker, "delivery_assignment")
	orderService := order_fulfillment.NewOrderService(deliveryPartnerService, deliveryAssignmentDAO, deliveryPartnerDAO, deliveryAssignmentProducer)

	// Initialize Kafka orderConsumer
	orderConsumer, err := consumer.NewOrderConsumer(broker, groupID, orderService)
	if err != nil {
		log.Fatalf("Error creating Kafka orderConsumer: %v", err)
	}

	// Start consuming messages
	orderConsumer.Start(topic)
}
