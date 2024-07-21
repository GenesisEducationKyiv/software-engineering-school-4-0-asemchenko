package util

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/messaging"
	"log"
)

func CloseRabbitMQ(consumer *messaging.RabbitMQConsumer) {
	err := consumer.Close()
	if err != nil {
		log.Fatalf("Failed to close RabbitMQ connection: %v", err)
	}
}
