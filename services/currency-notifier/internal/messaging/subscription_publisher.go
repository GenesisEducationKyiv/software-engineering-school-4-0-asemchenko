package messaging

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/api"
	"github.com/google/uuid"
)

type SubscriptionPublisher interface {
	PublishAddSubscription(id uuid.UUID, email string) error
	PublishRemoveSubscription(id uuid.UUID) error
}

type subscriptionPublisher struct {
	rabbitMqPublisher *RabbitMQPublisher
}

func NewSubscriptionPublisher(publisher *RabbitMQPublisher) SubscriptionPublisher {
	return &subscriptionPublisher{
		rabbitMqPublisher: publisher,
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

	return p.rabbitMqPublisher.PublishMessage(messageBody)
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

	return p.rabbitMqPublisher.PublishMessage(messageBody)
}
