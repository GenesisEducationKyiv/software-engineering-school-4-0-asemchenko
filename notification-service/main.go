package main

import (
	"encoding/json"
	"log"
	"notification-service/api"
	"notification-service/internal/messaging"
	"notification-service/internal/service"
	"notification-service/internal/util"
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
	defer util.CloseRabbitMQ(consumer)

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
