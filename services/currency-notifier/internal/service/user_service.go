package service

import (
	"currency-notifier/internal/exceptions"
	"currency-notifier/internal/messaging"
	"currency-notifier/internal/models"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/util"
	"github.com/emicklei/go-restful/v3/log"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	Repo                       *repository.UserRepository
	SubscriptionPublisher      messaging.SubscriptionPublisher
	SubscriptionStatusConsumer messaging.SubscriptionStatusListener
}

func NewUserService(repo *repository.UserRepository, publisher messaging.SubscriptionPublisher, consumer messaging.SubscriptionStatusListener) *UserService {
	return &UserService{
		Repo:                       repo,
		SubscriptionPublisher:      publisher,
		SubscriptionStatusConsumer: consumer,
	}
}

func (s *UserService) Subscribe(email string) (*models.User, error) {
	exists, err := s.Repo.IsEmailExists(email)
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

	err = s.SubscriptionPublisher.PublishAddSubscription(user.ID, user.Email)
	if err != nil {
		log.Printf("Failed to publish subscription message. Rollbacking - Deleting user %s", user.ID.String())
		util.WarnOnError(s.Repo.DeleteUser(user.ID.String()))
	}
	log.Printf("Published subscription message for user %s", user.ID.String())

	subscriptionStatusChan := s.SubscriptionStatusConsumer.SubscribeOnUser(user.ID)
	defer s.SubscriptionStatusConsumer.CloseSubscription(subscriptionStatusChan)

	select {
	case status := <-subscriptionStatusChan.Channel:
		log.Printf("Got subscription status: %v", status)
		if !status.IsSuccessfullyCreated() {
			log.Printf("Failed to create subscription: %v. Rollbacking - Deleting user %s", status.Details, user.ID.String())
			util.WarnOnError(s.Repo.DeleteUser(user.ID.String()))
			return nil, exceptions.ErrSubscriptionFailed(status.Details)
		}
	case <-time.After(30 * time.Second):
		util.WarnOnError(s.SubscriptionPublisher.PublishRemoveSubscription(user.ID))
		util.WarnOnError(s.Repo.DeleteUser(user.ID.String()))
		return nil, exceptions.ErrSubscriptionTimeout
	}

	return user, err
}

func (s *UserService) Unsubscribe(userId uuid.UUID) error {
	util.WarnOnError(s.Repo.DeleteUser(userId.String()))
	return s.SubscriptionPublisher.PublishRemoveSubscription(userId)
}
