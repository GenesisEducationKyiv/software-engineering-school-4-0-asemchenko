package service

import (
	"github.com/google/uuid"
	"notification-service/internal/models"
	"notification-service/internal/repository"
	"time"
)

type SubscriptionService interface {
	CreateSubscription(userId uuid.UUID, email string) error
	DeactivateSubscription(userId uuid.UUID) error
	GetAllActiveSubscriptions() ([]models.Subscription, error)
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) SubscriptionService {
	return &service{repo: repo}
}

type service struct {
	repo *repository.SubscriptionRepository
}

func (s *service) CreateSubscription(userId uuid.UUID, email string) error {
	return s.repo.CreateSubscription(&models.Subscription{
		ID:        userId,
		Email:     email,
		CreatedAt: time.Now(),
		IsActive:  true,
	})
}

func (s *service) DeactivateSubscription(userId uuid.UUID) error {
	return s.repo.DeactivateSubscription(userId)
}

func (s *service) GetAllActiveSubscriptions() ([]models.Subscription, error) {
	return s.repo.GetAllActiveSubscriptions()
}
