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

func (c *RateUpdateConsumer) Consume(message string) {
	var msg api.RateUpdateMessage
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	for s := range msg.Currencies {
		c.rateStorage.SaveRate(s, msg.Currencies[s])
	}
}
