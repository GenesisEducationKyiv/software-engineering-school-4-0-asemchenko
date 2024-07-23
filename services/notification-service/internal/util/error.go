package util

import (
	"go.uber.org/zap"
	"notification-service/internal/messaging"
)

func CloseRabbitMQ(consumer *messaging.RabbitMqClient) {
	err := consumer.Close()
	if err != nil {
		zap.S().Warnf("Failed to close RabbitMQ connection: %v", err)
	}
}

func LogError(err error) {
	if err != nil {
		zap.S().Errorf("Error occurred: %v", err)
	}
}
