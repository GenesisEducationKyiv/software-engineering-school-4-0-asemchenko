package service

import (
	"currency-notifier/internal/messaging"
	"encoding/json"
	"fmt"
	"log"
)

type EmailService struct {
	rabbitMQPublisher *messaging.RabbitMQPublisher
}

func NewEmailService(rabbitMQPublisher *messaging.RabbitMQPublisher) *EmailService {
	return &EmailService{
		rabbitMQPublisher: rabbitMQPublisher,
	}
}

func (s *EmailService) SendCurrencyRateEmail(to string, rate float64) error {
	plainTextContent := fmt.Sprintf("The current USD to UAH rate is %.2f.", rate)
	htmlContent := fmt.Sprintf("<strong>The current USD to UAH rate is %.2f.</strong>", rate)

	emailMessage := messaging.Message{
		TargetEmail:  to,
		Subject:      "Currency Rate Alert",
		MessagePlain: plainTextContent,
		MessageHtml:  htmlContent,
	}

	messageBody, err := json.Marshal(emailMessage)
	if err != nil {
		return err
	}

	err = s.rabbitMQPublisher.PublishMessage(messageBody)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	return nil
}
