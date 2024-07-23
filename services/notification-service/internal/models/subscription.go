package models

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}
