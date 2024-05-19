package models

import "time"

type Subscription struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	SubscribedAt time.Time `json:"subscribed_at" db:"subscribed_at"`
}
