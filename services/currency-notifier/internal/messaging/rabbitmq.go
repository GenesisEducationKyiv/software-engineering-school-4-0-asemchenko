package messaging

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMqClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	Url        string
}

func NewRabbitMqClient(url string) *RabbitMqClient {
	return &RabbitMqClient{
		Url: url,
	}
}

func (c *RabbitMqClient) Connect() {
	var err error
	c.connection, err = amqp.Dial(c.Url)
	if err != nil {
		zap.S().Fatalf("Failed to connect to RabbitMQ %v", err)
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		zap.S().Fatalf("Failed to open channel at RabbitMQ connection %v", err)
	}
}

func (c *RabbitMqClient) DeclareQueue(queueName string) {
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.S().Fatalf("Failed to declare a queue %v", err)
	}
}

func (c *RabbitMqClient) Close() {
	if c.channel != nil {
		err := c.channel.Close()
		if err != nil {
			zap.S().Fatalf("Failed to close RabbitMQ channel %v", err)
		}
	}
	if c.connection != nil {
		err := c.connection.Close()
		if err != nil {
			zap.S().Fatalf("Failed to close RabbitMQ connection %v", err)
		}
	}
}

func (c *RabbitMqClient) PublishMessage(messageBody []byte, queueName string) error {
	err := c.channel.Publish(
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
		zap.S().Warnf("Failed to publish message: %v", err)
		return err
	}

	return nil
}

func (c *RabbitMqClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	messagesChannel, err := c.channel.Consume(
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
