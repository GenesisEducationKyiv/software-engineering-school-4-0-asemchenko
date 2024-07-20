package service

import (
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/messaging"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

type EmailService struct {
	Dialer Dialer
	From   string
}

type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

func NewWithDialer(dialer Dialer, from string) *EmailService {
	return &EmailService{
		Dialer: dialer,
		From:   from,
	}
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM_EMAIL")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		portInt = 587 // default SMTP port
	}

	dialer := gomail.NewDialer(host, portInt, username, password)

	return &EmailService{
		Dialer: dialer,
		From:   from,
	}
}

func (s *EmailService) SendEmail(message messaging.Message) error {
	return s.sendEmail(message.TargetEmail, message.Subject, message.MessagePlain, message.MessageHtml)
}

func (s *EmailService) sendEmail(to string, subject string, plainTextContent string, htmlContent string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", plainTextContent)
	m.AddAlternative("text/html", htmlContent)

	if err := s.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email to %s: %w", to, err)
	}
	log.Printf("Email sent to %s", to)
	return nil
}
