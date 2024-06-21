package context

import (
	"currency-notifier/internal/jobs"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/service"
	"currency-notifier/internal/service/rate_provider/monobank"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

type AppContext struct {
	db           *Db
	cronInstance *cron.Cron

	SubscriptionRepo *repository.SubscriptionRepository
	RateRepository   *repository.ExchangeRateRepository

	MonobankRateProvider service.RateProvider

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
	ctx.db = NewDb(getDbConfig())
	ctx.db.Init()

	ctx.SubscriptionRepo = repository.NewSubscriptionRepository(ctx.db.Get())
	ctx.RateRepository = repository.NewExchangeRateRepository(ctx.db.Get())

	ctx.MonobankRateProvider = monobank.NewMonobankRateProvider(getEnv("MONOBANK_HOST_URL", "https://api.monobank.ua/bank/currency"))

	ctx.EmailService = service.NewEmailService()
	ctx.SubscriptionService = service.NewSubscriptionService(ctx.SubscriptionRepo)
	ctx.CurrencyService = service.NewCurrencyService(ctx.RateRepository, ctx.MonobankRateProvider)

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
	ctx.db.Close()
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

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getDbConfig() DbConfig {
	return DbConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "currency_notifier"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}
