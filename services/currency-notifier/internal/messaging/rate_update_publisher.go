package messaging

import (
	"currency-notifier/internal/shared/api"
	"encoding/json"
	"log"
)

type RatePublisher struct {
	rabbitMqPublisher *RabbitMQPublisher
}

func NewRatePublisher(publisher *RabbitMQPublisher) *RatePublisher {
	return &RatePublisher{
		rabbitMqPublisher: publisher,
	}
}

func (p *RatePublisher) OnCurrencyRateUpdated(usdUahRate float64) error {
	log.Printf("Sending rate update message: %v", usdUahRate)
	rateUpdateMessage := api.RateUpdateMessage{
		Currencies: map[string]float64{
			"usd_uah": usdUahRate,
		},
	}

	rateUpdateMessageBody, err := json.Marshal(rateUpdateMessage)
	if err != nil {
		return err
	}

	message := api.Message{
		Type: api.MessageTypeRateUpdate,
		Body: string(rateUpdateMessageBody),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.rabbitMqPublisher.PublishMessage(messageBody)
}
