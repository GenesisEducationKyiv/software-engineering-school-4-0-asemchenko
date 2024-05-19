package jobs

import (
	"currency-notifier/internal/service"
	"log"
)

func SendEmailsJob(currencyService *service.CurrencyService, subscriptionService *service.SubscriptionService, emailService *service.EmailService) {
	log.Printf("Sending emails with latest exchange rate")
	rate, err := currencyService.GetUSDtoUAHRate()
	if err != nil {
		log.Printf("Error fetching latest exchange rate: %v", err)
		return
	}
	log.Printf("Latest exchange rate: %f", rate)

	subscriptions, err := subscriptionService.GetAllSubscriptions()
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		return
	}
	log.Printf("Fetched %d subscriptions", len(subscriptions))

	for _, subscription := range subscriptions {
		err = emailService.SendCurrencyRateEmail(subscription.Email, rate)
		if err != nil {
			log.Printf("Error sending email to %s: %v", subscription.Email, err)
		} else {
			log.Printf("Email sent to %s", subscription.Email)
		}
	}
	log.Printf("Emails sent")
}
