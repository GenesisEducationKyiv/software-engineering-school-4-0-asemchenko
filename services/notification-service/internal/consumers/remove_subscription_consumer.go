package consumers

import (
	"encoding/json"
	"log"
	"notification-service/internal/service"
	"notification-service/internal/shared/api"
)

type RemoveSubscriptionConsumer interface {
	Consume(message string)
}

func NewRemoveSubscriptionConsumer(subscriptionService service.SubscriptionService) RemoveSubscriptionConsumer {
	return &removeSubscriptionConsumer{
		subscriptionService: subscriptionService,
	}
}

type removeSubscriptionConsumer struct {
	subscriptionService service.SubscriptionService
}

func (c *removeSubscriptionConsumer) Consume(message string) {
	var msg api.RemoveSubscriptionMessage
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	log.Printf("Received remove subscription message: %v", msg)
	err = c.subscriptionService.DeactivateSubscription(msg.UserID)
	if err != nil {
		log.Printf("Failed to deactivate subscription: %v", err)
	}
}
