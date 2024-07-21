package context

import (
	"currency-notifier/internal/jobs"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/service"
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/service/rate_provider/monobank"
	"currency-notifier/internal/service/rate_provider/nbu"
	"currency-notifier/internal/service/rate_provider/privatbank"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/messaging"
	sharedutil "github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/util"
	"github.com/robfig/cron/v3"
	"log"
)

type AppContext struct {
	db           *Db
	cronInstance *cron.Cron

	SubscriptionRepo *repository.SubscriptionRepository
	RateRepository   *repository.ExchangeRateRepository
	MessagePublisher *messaging.RabbitMQPublisher

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
	ctx.MessagePublisher = messaging.NewRabbitMQPublisher(
		sharedutil.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		sharedutil.GetEnv("RABBITMQ_QUEUE_NAME", "email"),
	)
	ctx.MessagePublisher.Connect()
	log.Printf("Connected to RabbitMQ")

	ctx.MonobankRateProvider = monobank.NewMonobankRateProvider(sharedutil.GetEnv("MONOBANK_HOST_URL", "https://api.monobank.ua"))
	ctx.PrivatBankRateProvider = privatbank.NewPrivatBankRateProvider(sharedutil.GetEnv("PRIVAT_BANK_HOST_URL", "https://api.privatbank.ua"))
	ctx.NbuRateProvider = nbu.NewNbuRateProvider(sharedutil.GetEnv("NBU_HOST_URL", "https://bank.gov.ua"))

	ctx.AggregatedRateProvider = rate_provider.NewAggregatedRateProvider(
		ctx.MonobankRateProvider,
		ctx.PrivatBankRateProvider,
		ctx.NbuRateProvider,
	)

	ctx.EmailService = service.NewEmailService(ctx.MessagePublisher)
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
		Host:     sharedutil.GetEnv("DB_HOST", "localhost"),
		Port:     sharedutil.GetEnv("DB_PORT", "5432"),
		User:     sharedutil.GetEnv("DB_USER", "postgres"),
		Password: sharedutil.GetEnv("DB_PASSWORD", "postgres"),
		Name:     sharedutil.GetEnv("DB_NAME", "currency_notifier"),
		SSLMode:  sharedutil.GetEnv("DB_SSLMODE", "disable"),
	}
}
