package util

import (
	"log"
	"notification-service/internal/messaging"
)

func CloseRabbitMQ(consumer *messaging.RabbitMQConsumer) {
	err := consumer.Close()
	if err != nil {
		log.Fatalf("Failed to close RabbitMQ connection: %v", err)
	}
}
