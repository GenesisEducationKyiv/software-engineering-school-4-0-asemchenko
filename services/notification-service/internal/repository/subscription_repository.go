package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"notification-service/internal/models"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
	query := "INSERT INTO subscriptions (id, email, created_at, is_active) VALUES ($1, $2, $3, $4)"
	_, err := r.DB.Exec(query, subscription.ID, subscription.Email, subscription.CreatedAt, subscription.IsActive)
	return err
}

func (r *SubscriptionRepository) GetAllActiveSubscriptions() ([]models.Subscription, error) {
	rows, err := r.DB.Query("SELECT id, email, created_at, is_active FROM subscriptions WHERE is_active=true")
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
		if err = rows.Scan(&subscription.ID, &subscription.Email, &subscription.CreatedAt, &subscription.IsActive); err != nil {
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

func (r *SubscriptionRepository) RemoveSubscription(userId uuid.UUID) error {
	query := "DELETE FROM subscriptions WHERE id=$1"
	_, err := r.DB.Exec(query, userId.String())
	return err
}
