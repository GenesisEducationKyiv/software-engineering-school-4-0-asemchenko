package consumers

import (
	"encoding/json"
	"go.uber.org/zap"
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
		zap.S().Warnf("Failed to parse message: %v", err)
		return
	}

	zap.S().Debugf("Received remove subscription message: %v", msg)
	err = c.subscriptionService.RemoveSubscription(msg.UserID)
	if err != nil {
		zap.S().Warnf("Failed to deactivate subscription: %v", err)
	}
}
