package service

import (
	"currency-notifier/internal/exceptions"
	"currency-notifier/internal/models"
	"currency-notifier/internal/repository"
	"time"
)

type SubscriptionService struct {
	Repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{Repo: repo}
}

func (s *SubscriptionService) Subscribe(email string) error {
	exists, err := s.Repo.SubscriptionExists(email)
	if err != nil {
		return err
	}
	if exists {
		return exceptions.ErrEmailAlreadySubscribed
	}

	subscription := &models.Subscription{
		Email:        email,
		SubscribedAt: time.Now(),
	}

	return s.Repo.AddSubscription(subscription)
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.Repo.GetAllSubscriptions()
}
