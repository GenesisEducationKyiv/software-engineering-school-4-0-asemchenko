package messaging

import (
	"currency-notifier/internal/shared/api"
	"encoding/json"
	"log"
)

type ConsumerController interface {
	StartConsuming(queueName string)
}

type consumerController struct {
	rabbitMqClient             *RabbitMqClient
	subscriptionStatusConsumer SubscriptionStatusListener
}

func NewConsumerController(
	rabbitMqClient *RabbitMqClient,
	subscriptionStatusConsumer SubscriptionStatusListener,
) ConsumerController {
	return &consumerController{
		rabbitMqClient:             rabbitMqClient,
		subscriptionStatusConsumer: subscriptionStatusConsumer,
	}
}

func (c *consumerController) StartConsuming(queueName string) {
	c.rabbitMqClient.DeclareQueue(queueName)
	messagesChannel, err := c.rabbitMqClient.Consume(queueName)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range messagesChannel {
		c.Consume(msg.Body)
	}
}

func (c *consumerController) Consume(message []byte) {
	// unpack message into api.Message type
	var msg api.Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	switch msg.Type {
	case api.MessageTypeSubscriptionStatus:
		c.subscriptionStatusConsumer.Consume(msg.Body)
	}
}
