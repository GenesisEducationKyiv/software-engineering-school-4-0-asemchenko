package context

import (
	"currency-notifier/internal/jobs"
	"currency-notifier/internal/messaging"
	"currency-notifier/internal/repository"
	"currency-notifier/internal/service"
	"currency-notifier/internal/service/rate_provider"
	"currency-notifier/internal/service/rate_provider/monobank"
	"currency-notifier/internal/service/rate_provider/nbu"
	"currency-notifier/internal/service/rate_provider/privatbank"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/util"
	"github.com/robfig/cron/v3"
	"log"
)

type AppContext struct {
	db           *db.Db
	cronInstance *cron.Cron

	UserRepository *repository.UserRepository
	RateRepository *repository.ExchangeRateRepository

	MessagePublisher      *messaging.RabbitMQPublisher
	RateUpdatePublisher   *messaging.RatePublisher
	SubscriptionPublisher messaging.SubscriptionPublisher

	MonobankRateProvider   rate_provider.RateProvider
	PrivatBankRateProvider rate_provider.RateProvider
	NbuRateProvider        rate_provider.RateProvider

	AggregatedRateProvider rate_provider.RateProvider

	UserService     *service.UserService
	CurrencyService *service.CurrencyService

	UpdateRateJob *jobs.UpdateRateJob
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

func (ctx *AppContext) Init() {
	ctx.db = db.NewDb(getDbConfig())
	ctx.db.Init()

	ctx.UserRepository = repository.NewUserRepository(ctx.db.Get())
	ctx.RateRepository = repository.NewExchangeRateRepository(ctx.db.Get())
	ctx.MessagePublisher = messaging.NewRabbitMQPublisher(
		util.GetEnv("RABBITMQ_URL", "amqp://username:password@localhost:5672/"),
		util.GetEnv("RABBITMQ_QUEUE_NAME", "email"),
	)
	ctx.MessagePublisher.Connect()
	log.Printf("Connected to RabbitMQ")
	ctx.RateUpdatePublisher = messaging.NewRatePublisher(ctx.MessagePublisher)
	ctx.SubscriptionPublisher = messaging.NewSubscriptionPublisher(ctx.MessagePublisher)

	ctx.MonobankRateProvider = monobank.NewMonobankRateProvider(util.GetEnv("MONOBANK_HOST_URL", "https://api.monobank.ua"))
	ctx.PrivatBankRateProvider = privatbank.NewPrivatBankRateProvider(util.GetEnv("PRIVAT_BANK_HOST_URL", "https://api.privatbank.ua"))
	ctx.NbuRateProvider = nbu.NewNbuRateProvider(util.GetEnv("NBU_HOST_URL", "https://bank.gov.ua"))

	ctx.AggregatedRateProvider = rate_provider.NewAggregatedRateProvider(
		ctx.MonobankRateProvider,
		ctx.PrivatBankRateProvider,
		ctx.NbuRateProvider,
	)

	ctx.UserService = service.NewUserService(ctx.UserRepository, ctx.SubscriptionPublisher)
	ctx.CurrencyService = service.NewCurrencyService(ctx.RateRepository, ctx.AggregatedRateProvider)
	ctx.CurrencyService.AddSubscriber(ctx.RateUpdatePublisher)

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

	ctx.cronInstance.Start()
}

func getDbConfig() db.DbConfig {
	return db.DbConfig{
		Host:     util.GetEnv("DB_HOST", "localhost"),
		Port:     util.GetEnv("DB_PORT", "5432"),
		User:     util.GetEnv("DB_USER", "postgres"),
		Password: util.GetEnv("DB_PASSWORD", "postgres"),
		Name:     util.GetEnv("DB_NAME", "currency_notifier"),
		SSLMode:  util.GetEnv("DB_SSLMODE", "disable"),
	}
}
