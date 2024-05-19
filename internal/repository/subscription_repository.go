package repository

import (
	"currency-notifier/internal/models"
	"database/sql"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) AddSubscription(subscription *models.Subscription) error {
	query := "INSERT INTO subscriptions (email, subscribed_at) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, subscription.Email, subscription.SubscribedAt)
	return err
}

func (r *SubscriptionRepository) SubscriptionExists(email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM subscriptions WHERE email=$1)"
	var exists bool
	err := r.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}
