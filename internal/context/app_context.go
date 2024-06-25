package context

import (
	"currency-notifier/internal/jobs"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/service"
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/service/rate_provider/monobank"
	"currency-notifier/internal/service/rate_provider/nbu"
	"currency-notifier/internal/service/rate_provider/privatbank"
	"currency-notifier/internal/util"
	"github.com/robfig/cron/v3"
	"log"
)

type AppContext struct {
	db           *Db
	cronInstance *cron.Cron

	SubscriptionRepo *repository.SubscriptionRepository
	RateRepository   *repository.ExchangeRateRepository

	MonobankRateProvider   rate_provider.RateProvider
	PrivatBankRateProvider rate_provider.RateProvider
	NbuRateProvider        rate_provider.RateProvider

	AggregatedRateProvider rate_provider.RateProvider

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

	ctx.MonobankRateProvider = monobank.NewMonobankRateProvider(util.GetEnv("MONOBANK_HOST_URL", "https://api.monobank.ua"))
	ctx.PrivatBankRateProvider = privatbank.NewPrivatBankRateProvider(util.GetEnv("PRIVAT_BANK_HOST_URL", "https://api.privatbank.ua"))
	ctx.NbuRateProvider = nbu.NewNbuRateProvider(util.GetEnv("NBU_HOST_URL", "https://bank.gov.ua"))

	ctx.AggregatedRateProvider = rate_provider.NewAggregatedRateProvider(
		ctx.MonobankRateProvider,
		ctx.PrivatBankRateProvider,
		ctx.NbuRateProvider,
	)

	ctx.EmailService = service.NewEmailService()
	ctx.SubscriptionService = service.NewSubscriptionService(ctx.SubscriptionRepo)
	ctx.CurrencyService = service.NewCurrencyService(ctx.RateRepository, ctx.AggregatedRateProvider)

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

func getDbConfig() DbConfig {
	return DbConfig{
		Host:     util.GetEnv("DB_HOST", "localhost"),
		Port:     util.GetEnv("DB_PORT", "5432"),
		User:     util.GetEnv("DB_USER", "postgres"),
		Password: util.GetEnv("DB_PASSWORD", "postgres"),
		Name:     util.GetEnv("DB_NAME", "currency_notifier"),
		SSLMode:  util.GetEnv("DB_SSLMODE", "disable"),
	}
}
