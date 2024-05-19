package service

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

type EmailService struct {
	Dialer *gomail.Dialer
	From   string
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

func (s *EmailService) SendCurrencyRateEmail(to string, rate float64) error {
	subject := "Currency Rate Alert"
	plainTextContent := fmt.Sprintf("The current USD to UAH rate is %.2f.", rate)
	htmlContent := fmt.Sprintf("<strong>The current USD to UAH rate is %.2f.</strong>", rate)

	return s.SendEmail(to, subject, plainTextContent, htmlContent)
}

func (s *EmailService) SendEmail(to string, subject string, plainTextContent string, htmlContent string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", plainTextContent)
	m.AddAlternative("text/html", htmlContent)

	if err := s.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email to %s: %w", to, err)
	}
	return nil
}
