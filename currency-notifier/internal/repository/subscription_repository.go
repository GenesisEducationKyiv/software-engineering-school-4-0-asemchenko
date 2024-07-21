package repository

import (
	"currency-notifier/internal/models"
	"database/sql"
	"log"
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

func (r *SubscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	rows, err := r.DB.Query("SELECT email, subscribed_at FROM subscriptions")
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var subscriptions []models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		if err = rows.Scan(&subscription.Email, &subscription.SubscribedAt); err != nil {
			log.Printf("Error scanning subscription: %v", err)
			continue
		}
		subscriptions = append(subscriptions, subscription)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error with rows in GetAllSubscriptions: %v", err)
		return nil, err
	}
	return subscriptions, nil
}
