package consumers

import (
	"encoding/json"
	"log"
	"notification-service/internal/service"
	"notification-service/internal/shared/api"
)

type RateUpdateConsumer struct {
	rateStorage service.RateStorage
}

func NewRateUpdateConsumer(rateStorage service.RateStorage) *RateUpdateConsumer {
	return &RateUpdateConsumer{
		rateStorage: rateStorage,
	}
}

func (c *RateUpdateConsumer) Consume(message string) {
	var msg api.RateUpdateMessage
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	log.Printf("Received rate update message: %v", msg.Currencies)
	for s := range msg.Currencies {
		c.rateStorage.SaveRate(s, msg.Currencies[s])
	}
}
