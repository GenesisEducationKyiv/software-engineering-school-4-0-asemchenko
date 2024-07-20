package messaging

import (
	"github.com/streadway/amqp"
)

type RabbitMQCredentials struct {
	Url       string
	QueueName string
}

type RabbitMQConsumer struct {
	Conn        *amqp.Connection
	Channel     *amqp.Channel
	Credentials RabbitMQCredentials
}

func NewRabbitMQConsumer(credentials RabbitMQCredentials) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		Credentials: credentials,
	}
}

func (r *RabbitMQConsumer) Connect() error {
	conn, err := amqp.Dial(r.Credentials.Url)
	if err != nil {
		return err
	}
	r.Conn = conn

	ch, err := r.Conn.Channel()
	if err != nil {
		return err
	}
	r.Channel = ch

	err = r.DeclareQueue()
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) Close() error {
	if r.Channel != nil {
		err := r.Channel.Close()
		if err != nil {
			return err
		}
	}
	return r.Conn.Close()
}

func (r *RabbitMQConsumer) DeclareQueue() error {
	_, err := r.Channel.QueueDeclare(
		r.Credentials.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMQConsumer) Consume() (<-chan amqp.Delivery, error) {
	messagesChannel, err := r.Channel.Consume(
		r.Credentials.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messagesChannel, nil
}
