package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQPublisher struct {
	Url        string
	QueueName  string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQPublisher(url string, queueName string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		Url:       url,
		QueueName: queueName,
	}
}

func (p *RabbitMQPublisher) Connect() {
	var err error
	p.connection, err = amqp.Dial(p.Url)
	if err != nil {
		log.Fatal(err)
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// Declare a queue
	_, err = p.channel.QueueDeclare(
		p.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}
}

func (p *RabbitMQPublisher) Close() {
	if p.channel != nil {
		err := p.channel.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	if p.connection != nil {
		err := p.connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (p *RabbitMQPublisher) PublishMessage(messageBody []byte) error {
	err := p.channel.Publish(
		"",
		p.QueueName,
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
