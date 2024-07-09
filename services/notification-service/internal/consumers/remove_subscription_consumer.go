package consumers

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
	"log"
	"notification-service/internal/service"
)

type RemoveSubscriptionConsumer interface {
	Consume(message []byte)
}

func NewRemoveSubscriptionConsumer(subscriptionService service.SubscriptionService) RemoveSubscriptionConsumer {
	return &removeSubscriptionConsumer{
		subscriptionService: subscriptionService,
	}
}

type removeSubscriptionConsumer struct {
	subscriptionService service.SubscriptionService
}

func (c *removeSubscriptionConsumer) Consume(message []byte) {
	var msg api.RemoveSubscriptionMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	err = c.subscriptionService.DeactivateSubscription(msg.UserID)
	if err != nil {
		log.Printf("Failed to deactivate subscription: %v", err)
	}
}
