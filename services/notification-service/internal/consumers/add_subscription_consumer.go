package consumers

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
	"log"
	"notification-service/internal/service"
)

type AddSubscriptionConsumer interface {
	Consume(message []byte)
}

func NewAddSubscriptionConsumer(subscriptionService service.SubscriptionService) AddSubscriptionConsumer {
	return &addSubscriptionConsumer{
		subscriptionService: subscriptionService,
	}
}

type addSubscriptionConsumer struct {
	subscriptionService service.SubscriptionService
}

func (c *addSubscriptionConsumer) Consume(message []byte) {
	var msg api.AddSubscriptionMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	err = c.subscriptionService.CreateSubscription(msg.UserID, msg.Email)
	if err != nil {
		log.Printf("Failed to create subscription: %v", err)
	}
}
