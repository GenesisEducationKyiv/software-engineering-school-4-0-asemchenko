package main

import (
	"currency-notifier/internal/controller"
	"currency-notifier/internal/jobs"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/service"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"

	_ "currency-notifier/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB
var cronInstance *cron.Cron

// @title UAH currency application
// @version 1.0
// @description API for current USD-UAH exchange rate and for email-subscribing on the currency rate
// @host localhost:8080
// @basePath /api

//go:generate go run github.com/swaggo/swag/cmd/swag init

func main() {
	r := mux.NewRouter()

	err := initDb()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	rateRepository := repository.NewExchangeRateRepository(db)

	emailService := service.NewEmailService()
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	currencyService := service.NewCurrencyService(rateRepository)

	err = currencyService.Init()
	if err != nil {
		log.Fatal("Error initializing currency service: ", err)
	}

	subscriptionController := controller.NewSubscriptionController(subscriptionService)
	rateController := controller.NewRateController(currencyService)

	r.HandleFunc("/api/rate", rateController.GetRate).Methods("GET")
	r.HandleFunc("/api/subscribe", subscriptionController.Subscribe).Methods("POST")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	initCron(emailService, currencyService, subscriptionService)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initCron(emailService *service.EmailService,
	currencyService *service.CurrencyService,
	subscriptionService *service.SubscriptionService) {

	cronInstance = cron.New()

	_, err := cronInstance.AddFunc("@every 1m", func() {
		jobs.UpdateExchangeRateJob(currencyService)
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = cronInstance.AddFunc("@every 1m", func() {
		jobs.SendEmailsJob(currencyService, subscriptionService, emailService)
	})
	if err != nil {
		log.Fatal()
	}

	cronInstance.Start()
}

func initDb() error {
	// If no env variables - assume it's 'local' environment and use default values

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "currency_notifier")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = runMigrations(db)

	return err
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not start migration: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %w", err)
	}
	log.Println("Migrations ran successfully")
	return nil
}
