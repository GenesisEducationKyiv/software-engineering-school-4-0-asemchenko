package service

import (
	"currency-notifier/internal/exceptions"
	"currency-notifier/internal/messaging"
	"currency-notifier/internal/models"
	"currency-notifier/internal/repository"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	Repo           *repository.UserRepository
	eventPublisher messaging.SubscriptionPublisher
}

func NewUserService(repo *repository.UserRepository, publisher messaging.SubscriptionPublisher) *UserService {
	return &UserService{Repo: repo, eventPublisher: publisher}
}

func (s *UserService) Subscribe(email string) (*models.User, error) {
	exists, err := s.Repo.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, exceptions.ErrEmailAlreadyExists
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now(),
	}

	err = s.Repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	err = s.eventPublisher.PublishAddSubscription(user.ID, user.Email)
	return user, err
}

func (s *UserService) Unsubscribe(userId uuid.UUID) error {
	return s.eventPublisher.PublishRemoveSubscription(userId)
}
