package consumers

import (
	"encoding/json"
	"go.uber.org/zap"
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
		zap.S().Warnf("Failed to parse message: %v", err)
		return
	}

	zap.S().Debugf("Received add subscription message: %v", msg)
	err = c.subscriptionService.CreateSubscription(msg.UserID, msg.Email)
	if err != nil {
		util.LogError(c.subscriptionStatusPublisher.PublishSubscriptionStatus(msg.UserID, api.SubscriptionStatusFailedToCreate))
		zap.S().Warnf("Sending failed subscription status message for user %v", msg.UserID)
		zap.S().Warnf("Failed to create subscription: %v", err)
		return
	}
	util.LogError(c.subscriptionStatusPublisher.PublishSubscriptionStatus(msg.UserID, api.SubscriptionStatusSuccessfullyCreated))
	zap.S().Debugf("Sending success subscription status message for user %v", msg.UserID)
}
