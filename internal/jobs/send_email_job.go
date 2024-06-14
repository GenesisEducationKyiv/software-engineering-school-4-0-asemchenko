package jobs

import (
	"currency-notifier/internal/models"
	"currency-notifier/internal/service"
	"log"
)

type SendEmailJob struct {
	currencyService     *service.CurrencyService
	subscriptionService *service.SubscriptionService
	emailService        *service.EmailService
}

func NewEmailSender(currencyService *service.CurrencyService, subscriptionService *service.SubscriptionService, emailService *service.EmailService) *SendEmailJob {
	return &SendEmailJob{
		currencyService:     currencyService,
		subscriptionService: subscriptionService,
		emailService:        emailService,
	}
}

func (j SendEmailJob) SendEmails() {
	log.Printf("Sending emails with latest exchange rate")
	rate, failed := j.getCurrentExchangeRate()
	if failed {
		return
	}

	subscriptions, failed := j.getActiveSubscriptions()
	if failed {
		return
	}

	j.sendRatesToActiveSubscribers(&subscriptions, rate)
}

func (j SendEmailJob) getCurrentExchangeRate() (float64, bool) {
	rate, err := j.currencyService.GetUSDtoUAHRate()
	if err != nil {
		log.Printf("Error fetching latest exchange rate: %v", err)
		return 0, true
	}
	log.Printf("Latest exchange rate: %f", rate)

	return rate, false
}

func (j SendEmailJob) getActiveSubscriptions() ([]models.Subscription, bool) {
	subscriptions, err := j.subscriptionService.GetAllSubscriptions()
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		return nil, true
	}
	log.Printf("Fetched %d subscriptions", len(subscriptions))

	return subscriptions, false
}

func (j SendEmailJob) sendRatesToActiveSubscribers(subscriptions *[]models.Subscription, rate float64) {
	for _, subscription := range *subscriptions {
		err := j.emailService.SendCurrencyRateEmail(subscription.Email, rate)
		if err != nil {
			log.Printf("Error sending email to %j: %v", subscription.Email, err)
		} else {
			log.Printf("Email sent to %j", subscription.Email)
		}
	}
	log.Printf("Emails sent")
}
