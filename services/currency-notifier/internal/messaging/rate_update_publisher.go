package messaging

import (
	"currency-notifier/internal/shared/api"
	"encoding/json"
	"go.uber.org/zap"
)

type RatePublisher struct {
	rabbitMqPublisher *RabbitMqClient
	queueName         string
}

func NewRatePublisher(publisher *RabbitMqClient, queueName string) *RatePublisher {
	publisher.DeclareQueue(queueName)
	return &RatePublisher{
		rabbitMqPublisher: publisher,
		queueName:         queueName,
	}
}

func (p *RatePublisher) PublishRateUpdated(usdUahRate float64) error {
	zap.S().Debugf("Sending rate update message: %v", usdUahRate)
	rateUpdateMessage := api.RateUpdatedMessage{
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

	return p.rabbitMqPublisher.PublishMessage(messageBody, p.queueName)
}
