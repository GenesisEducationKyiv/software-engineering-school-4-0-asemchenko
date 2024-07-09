package main

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/util"
	"log"
	"notification-service/internal/messaging"
	"notification-service/internal/service"
	rabbitUtils "notification-service/internal/util"
)

func main() {
	rabbitMQCredentials := messaging.RabbitMQCredentials{
		Url:       util.GetEnv("RABBITMQ_URL", "amqp://username:password@localhost:5672/"),
		QueueName: util.GetEnv("RABBITMQ_QUEUE_NAME", "email"),
	}
	consumer := messaging.NewRabbitMQConsumer(rabbitMQCredentials)
	err := consumer.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Printf("Connected to RabbitMQ")
	defer rabbitUtils.CloseRabbitMQ(consumer)

	emailService := service.NewEmailService()

	messagesChannel, err := consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range messagesChannel {
		var emailMessage api.Message
		err := json.Unmarshal(msg.Body, &emailMessage)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		err = emailService.SendEmail(emailMessage)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}
}
