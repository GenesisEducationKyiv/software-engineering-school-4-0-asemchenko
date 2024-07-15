package util

import (
	"log"
	"notification-service/internal/messaging"
)

func CloseRabbitMQ(consumer *messaging.RabbitMqClient) {
	err := consumer.Close()
	if err != nil {
		log.Fatalf("Failed to close RabbitMQ connection: %v", err)
	}
}

func WarnOnError(err error) {
	if err != nil {
		log.Printf("Error occurred: %v", err)
	}
}
