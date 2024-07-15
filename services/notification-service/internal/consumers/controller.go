package consumers

import (
	"encoding/json"
	"log"
	"notification-service/internal/messaging"
	"notification-service/internal/shared/api"
	"notification-service/internal/util"
)

type ConsumerController interface {
	StartConsuming()
}

type consumerController struct {
	rabbitMqConsumer           *messaging.RabbitMqClient
	rateUpdateConsumer         *RateUpdateConsumer
	addSubscriberConsumer      AddSubscriptionConsumer
	removeSubscriptionConsumer RemoveSubscriptionConsumer
	queueName                  string
}

func NewConsumerController(
	rabbitMqConsumer *messaging.RabbitMqClient,
	rateUpdateConsumer *RateUpdateConsumer,
	addSubscriberConsumer AddSubscriptionConsumer,
	removeSubscriptionConsumer RemoveSubscriptionConsumer,
	queueName string,
) ConsumerController {
	util.WarnOnError(rabbitMqConsumer.DeclareQueue(queueName))
	return &consumerController{
		rabbitMqConsumer:           rabbitMqConsumer,
		rateUpdateConsumer:         rateUpdateConsumer,
		addSubscriberConsumer:      addSubscriberConsumer,
		removeSubscriptionConsumer: removeSubscriptionConsumer,
		queueName:                  queueName,
	}
}

func (c *consumerController) StartConsuming() {
	messagesChannel, err := c.rabbitMqConsumer.Consume(c.queueName)
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
	case api.MessageTypeRateUpdate:
		c.rateUpdateConsumer.Consume(msg.Body)
	case api.MessageTypeAddSubscription:
		c.addSubscriberConsumer.Consume(msg.Body)
	case api.MessageTypeRemoveSubscription:
		c.removeSubscriptionConsumer.Consume(msg.Body)
	}
}
