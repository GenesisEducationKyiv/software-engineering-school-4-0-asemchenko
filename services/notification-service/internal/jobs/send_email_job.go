package jobs

import (
	"go.uber.org/zap"
	"notification-service/internal/models"
	"notification-service/internal/service"
)

type SendEmailJob struct {
	rateStorage         service.RateStorage
	subscriptionService service.SubscriptionService
	emailService        *service.EmailService
}

func NewSendEmailJob(rateStorage service.RateStorage, subscriptionService service.SubscriptionService, emailService *service.EmailService) *SendEmailJob {
	return &SendEmailJob{
		rateStorage:         rateStorage,
		subscriptionService: subscriptionService,
		emailService:        emailService,
	}
}

func (j *SendEmailJob) SendEmails() {
	zap.S().Debugf("Sending emails with latest exchange rate")
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

func (j *SendEmailJob) getCurrentExchangeRate() (float64, bool) {
	rate, err := j.rateStorage.GetLatestRate()
	if err != nil {
		zap.S().Errorf("Error fetching latest exchange rate: %v", err)
		return 0, true
	}
	zap.S().Debugf("Latest exchange rate: %f", rate)

	return rate, false
}

func (j *SendEmailJob) getActiveSubscriptions() ([]models.Subscription, bool) {
	subscriptions, err := j.subscriptionService.GetAllActiveSubscriptions()
	if err != nil {
		zap.S().Errorf("Error fetching subscriptions: %v", err)
		return nil, true
	}
	zap.S().Debugf("Fetched %d subscriptions", len(subscriptions))

	return subscriptions, false
}

func (j *SendEmailJob) sendRatesToActiveSubscribers(subscriptions *[]models.Subscription, rate float64) {
	for _, subscription := range *subscriptions {
		err := j.emailService.SendCurrencyRateEmail(subscription.Email, rate)
		if err != nil {
			zap.S().Errorf("Error sending email to %s: %v", subscription.Email, err)
		} else {
			zap.S().Debugf("Email for %s sent to processing", subscription.Email)
		}
	}
	zap.S().Debugf("Emails sent")
}
