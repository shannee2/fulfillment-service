package main

import (
	"context"
	"fulfillment/dao"
	"fulfillment/kafka/consumer"
	"fulfillment/kafka/producer"
	"fulfillment/service/delivery_partner"
	"fulfillment/service/order_fulfillment"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize Database
	db, err := dao.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Kafka configuration
	broker := "localhost:9092"
	topic := "orders"
	groupID := "order-consumer-group"

	// Initialize DAOs
	deliveryPartnerDAO := dao.NewDeliveryPartnerDAO(db)
	deliveryAssignmentDAO := dao.NewDeliveryAssignmentDAO(db)

	// Initialize Services
	deliveryPartnerService := delivery_partner.NewDeliveryPartnerService(deliveryPartnerDAO)
	producer, _ := producer.NewDeliveryAssignmentProducer(broker, "delivery_assignment")
	orderService := order_fulfillment.NewOrderService(deliveryPartnerService, deliveryAssignmentDAO, deliveryPartnerDAO, producer)

	// Initialize Kafka Consumer
	orderConsumer, err := consumer.NewOrderConsumer(broker, groupID, orderService)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer orderConsumer.Close() // Ensure consumer closes on exit

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Channel to listen for OS signals (CTRL+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start Kafka consumer in a goroutine
	go func() {
		orderConsumer.Start(ctx, topic)
	}()

	log.Println("Order Fulfillment Service is running...")

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down Order Fulfillment Service...")

	// Cleanup
	cancel()              // Cancel context to stop consumer
	orderConsumer.Close() // Close Kafka consumer properly
	log.Println("Service stopped gracefully")
}
