package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Email     string    `json:"email" db:"email"`
	ID        uuid.UUID `json:"id" db:"id"`
}
