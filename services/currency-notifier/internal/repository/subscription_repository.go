package repository

import (
	"currency-notifier/internal/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (id, email, created_at) VALUES ($1, $2, $3)"
	_, err := r.DB.Exec(query, user.ID, user.Email, user.CreatedAt)
	return err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)"
	var exists bool
	err := r.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}
