package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMqClient struct {
	Conn          *amqp.Connection
	Channel       *amqp.Channel
	ConnectionUrl string
}

func NewRabbitMQConsumer(connectionUrl string) *RabbitMqClient {
	return &RabbitMqClient{
		ConnectionUrl: connectionUrl,
	}
}

func (c *RabbitMqClient) Connect() error {
	conn, err := amqp.Dial(c.ConnectionUrl)
	if err != nil {
		return err
	}
	c.Conn = conn

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	c.Channel = ch

	return nil
}

func (c *RabbitMqClient) Close() error {
	if c.Channel != nil {
		err := c.Channel.Close()
		if err != nil {
			return err
		}
	}
	return c.Conn.Close()
}

func (c *RabbitMqClient) DeclareQueue(queueName string) error {
	_, err := c.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (c *RabbitMqClient) PublishMessage(messageBody []byte, queueName string) error {
	err := c.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	return nil
}

func (c *RabbitMqClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	messagesChannel, err := c.Channel.Consume(
		queueName,
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
