package messaging

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
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
