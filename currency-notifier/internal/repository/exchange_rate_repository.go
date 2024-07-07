package repository

import (
	"database/sql"
	"log"
)

type ExchangeRateRepository struct {
	DB *sql.DB
}

func NewExchangeRateRepository(db *sql.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{
		DB: db,
	}
}

func (repo *ExchangeRateRepository) SaveRate(rate float64) error {
	_, err := repo.DB.Exec("INSERT INTO exchange_rates (rate) VALUES ($1)", rate)
	if err != nil {
		log.Printf("Error saving exchange rate: %v", err)
		return err
	}
	return nil
}

func (repo *ExchangeRateRepository) GetLatestRate() (float64, error) {
	var rate float64
	err := repo.DB.QueryRow("SELECT rate FROM exchange_rates ORDER BY fetched_at DESC LIMIT 1").Scan(&rate)
	if err != nil {
		log.Printf("Error fetching latest exchange rate: %v", err)
		return 0, err
	}
	return rate, nil
}
