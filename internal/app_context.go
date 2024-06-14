package internal

import (
	"currency-notifier/internal/jobs"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/service"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

type AppContext struct {
	db           *sql.DB
	cronInstance *cron.Cron

	SubscriptionRepo *repository.SubscriptionRepository
	RateRepository   *repository.ExchangeRateRepository

	EmailService        *service.EmailService
	SubscriptionService *service.SubscriptionService
	CurrencyService     *service.CurrencyService

	SendEmailJob  *jobs.SendEmailJob
	UpdateRateJob *jobs.UpdateRateJob
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

func (ctx *AppContext) Init() {
	ctx.initDb()

	ctx.SubscriptionRepo = repository.NewSubscriptionRepository(ctx.db)
	ctx.RateRepository = repository.NewExchangeRateRepository(ctx.db)

	ctx.EmailService = service.NewEmailService()
	ctx.SubscriptionService = service.NewSubscriptionService(ctx.SubscriptionRepo)
	ctx.CurrencyService = service.NewCurrencyService(ctx.RateRepository)

	ctx.SendEmailJob = jobs.NewEmailSender(ctx.CurrencyService, ctx.SubscriptionService, ctx.EmailService)
	ctx.UpdateRateJob = jobs.NewUpdateRateJob(ctx.CurrencyService)

	err := ctx.CurrencyService.Init()
	if err != nil {
		log.Fatal("Error initializing currency service: ", err)
	}

	ctx.initCron()
}

func (ctx *AppContext) Close() {
	ctx.cronInstance.Stop()
	closeDb(ctx.db)
}

func (ctx *AppContext) initCron() {
	ctx.cronInstance = cron.New()

	_, err := ctx.cronInstance.AddFunc("@hourly", ctx.UpdateRateJob.UpdateExchangeRate)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ctx.cronInstance.AddFunc("@daily", ctx.SendEmailJob.SendEmails)
	if err != nil {
		log.Fatal()
	}

	ctx.cronInstance.Start()
}

func (ctx *AppContext) initDb() {
	// If no env variables - assume it's 'local' environment and use default values

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "currency_notifier")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	log.Printf("Connecting to database %s:%s as user %s to database %s\n", dbHost, dbPort, dbUser, dbName)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	var err error
	ctx.db, err = sql.Open("postgres", connStr)

	if err == nil {
		err = runMigrations(ctx.db)
	}

	if err != nil {
		log.Fatal("Failed to initialize DB", err)
	}
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

func closeDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
