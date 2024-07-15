package messaging

import (
	"currency-notifier/internal/shared/api"
	"encoding/json"
	"github.com/google/uuid"
)

type SubscriptionPublisher interface {
	PublishAddSubscription(id uuid.UUID, email string) error
	PublishRemoveSubscription(id uuid.UUID) error
}

type subscriptionPublisher struct {
	rabbitMqPublisher *RabbitMqClient
	queueName         string
}

func NewSubscriptionPublisher(publisher *RabbitMqClient, queueName string) SubscriptionPublisher {
	publisher.DeclareQueue(queueName)
	return &subscriptionPublisher{
		rabbitMqPublisher: publisher,
		queueName:         queueName,
	}
}

func (p *subscriptionPublisher) PublishAddSubscription(id uuid.UUID, email string) error {
	addMsg := api.AddSubscriptionMessage{
		UserID: id,
		Email:  email,
	}

	addMsgBody, err := json.Marshal(addMsg)
	if err != nil {
		return err
	}

	message := api.Message{
		Type: api.MessageTypeAddSubscription,
		Body: string(addMsgBody),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.rabbitMqPublisher.PublishMessage(messageBody, p.queueName)
}

func (p *subscriptionPublisher) PublishRemoveSubscription(id uuid.UUID) error {
	removeMsg := api.RemoveSubscriptionMessage{
		UserID: id,
	}

	removeMsgBody, err := json.Marshal(removeMsg)
	if err != nil {
		return err
	}

	message := api.Message{
		Type: api.MessageTypeRemoveSubscription,
		Body: string(removeMsgBody),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.rabbitMqPublisher.PublishMessage(messageBody, p.queueName)
}
