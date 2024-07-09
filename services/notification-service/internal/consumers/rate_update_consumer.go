package consumers

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
	"log"
	"notification-service/internal/service"
)

type RateUpdateConsumer struct {
	rateStorage service.RateStorage
}

func NewRateUpdateConsumer(rateStorage service.RateStorage) *RateUpdateConsumer {
	return &RateUpdateConsumer{
		rateStorage: rateStorage,
	}
}

func (c *RateUpdateConsumer) Consume(message []byte) {
	var msg api.RateUpdateMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	log.Printf("Received rate update message: %v", msg)
	for s := range msg.Currencies {
		c.rateStorage.SaveRate(s, msg.Currencies[s])
	}
}
