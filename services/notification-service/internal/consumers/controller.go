package consumers

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
	"log"
)

type ConsumerController struct {
	rateUpdateConsumer         *RateUpdateConsumer
	addSubscriberConsumer      *AddSubscriptionConsumer
	removeSubscriptionConsumer *RemoveSubscriptionConsumer
}

func (c *ConsumerController) Consume(message string) {
	// unpack message into api.Message type
	var msg api.Message
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	switch msg.Type {
	case api.MessageTypeRateUpdate:
		c.rateUpdateConsumer.Consume(message)
		break
	case api.MessageTypeAddSubscription:
		c.addSubscriberConsumer.Consume(message)
		break
	case api.MessageTypeRemoveSubscription:
		c.removeSubscriptionConsumer.Consume(message)
		break
	}
}
