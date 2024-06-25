package models

import "time"

type Subscription struct {
	SubscribedAt time.Time `json:"subscribed_at" db:"subscribed_at"`
	Email        string    `json:"email" db:"email"`
	ID           int       `json:"id" db:"id"`
}
