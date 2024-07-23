package consumers

import (
	"encoding/json"
	"log"
	"notification-service/internal/publishers"
	"notification-service/internal/service"
	"notification-service/internal/shared/api"
	"notification-service/internal/util"
)

type AddSubscriptionConsumer interface {
	Consume(message string)
}

func NewAddSubscriptionConsumer(
	subscriptionService service.SubscriptionService,
	subscriptionStatusPublisher publishers.SubscriptionStatusPublisher,
) AddSubscriptionConsumer {
	return &addSubscriptionConsumer{
		subscriptionService:         subscriptionService,
		subscriptionStatusPublisher: subscriptionStatusPublisher,
	}
}

type addSubscriptionConsumer struct {
	subscriptionService         service.SubscriptionService
	subscriptionStatusPublisher publishers.SubscriptionStatusPublisher
}

func (c *addSubscriptionConsumer) Consume(message string) {
	var msg api.AddSubscriptionMessage
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	log.Printf("Received add subscription message: %v", msg)
	err = c.subscriptionService.CreateSubscription(msg.UserID, msg.Email)
	if err != nil {
		util.WarnOnError(c.subscriptionStatusPublisher.PublishSubscriptionStatus(msg.UserID, api.SubscriptionStatusFailedToCreate))
		log.Printf("Sending failed subscription status message for user %v", msg.UserID)
		log.Printf("Failed to create subscription: %v", err)
		return
	}
	util.WarnOnError(c.subscriptionStatusPublisher.PublishSubscriptionStatus(msg.UserID, api.SubscriptionStatusSuccessfullyCreated))
	log.Printf("Sending success subscription status message for user %v", msg.UserID)
}
