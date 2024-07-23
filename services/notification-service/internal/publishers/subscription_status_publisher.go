package publishers

import (
	"encoding/json"
	"github.com/google/uuid"
	"notification-service/internal/messaging"
	"notification-service/internal/metrics"
	"notification-service/internal/shared/api"
	"notification-service/internal/util"
)

type SubscriptionStatusPublisher interface {
	PublishSubscriptionStatus(userId uuid.UUID, status string) error
}

type subscriptionPublisher struct {
	rabbitMqPublisher *messaging.RabbitMqClient
	queueName         string
}

func NewSubscriptionPublisher(publisher *messaging.RabbitMqClient, queueName string) SubscriptionStatusPublisher {
	util.LogError(publisher.DeclareQueue(queueName))
	return &subscriptionPublisher{
		rabbitMqPublisher: publisher,
		queueName:         queueName,
	}
}

func (p *subscriptionPublisher) PublishSubscriptionStatus(userId uuid.UUID, status string) error {
	statusMsg := api.SubscriptionStatusMessage{
		UserID: userId,
		Status: status,
	}

	statusMsgBody, err := json.Marshal(statusMsg)
	if err != nil {
		return err
	}

	message := api.Message{
		Type: api.MessageTypeSubscriptionStatus,
		Body: string(statusMsgBody),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.rabbitMqPublisher.PublishMessage(messageBody, p.queueName)
	if err != nil {
		return err
	}
	metrics.RabbitMQSentMessages.WithLabelValues(api.MessageTypeSubscriptionStatus).Inc()

	return nil
}
