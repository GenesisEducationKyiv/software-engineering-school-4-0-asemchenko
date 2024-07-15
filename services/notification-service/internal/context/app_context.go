package context

import (
	"github.com/robfig/cron/v3"
	"log"
	"notification-service/internal/consumers"
	"notification-service/internal/jobs"
	"notification-service/internal/messaging"
	"notification-service/internal/publishers"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/shared/db"
	"notification-service/internal/shared/util"
)

type AppContext struct {
	db *db.Db

	cronInstance *cron.Cron

	subscriptionRepository *repository.SubscriptionRepository

	emailService        *service.EmailService
	subscriptionService service.SubscriptionService

	rateStorage service.RateStorage

	rabbitMqConsumer            *messaging.RabbitMqClient
	rateUpdateConsumer          *consumers.RateUpdateConsumer
	addSubscriptionConsumer     consumers.AddSubscriptionConsumer
	removeSubscriptionConsumer  consumers.RemoveSubscriptionConsumer
	consumerController          consumers.ConsumerController
	subscriptionStatusPublisher publishers.SubscriptionStatusPublisher

	sendEmailJob *jobs.SendEmailJob
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

func (ctx *AppContext) Init() {
	ctx.db = db.NewDb(getDbConfig())
	ctx.db.Init()

	ctx.cronInstance = cron.New()

	ctx.subscriptionRepository = repository.NewSubscriptionRepository(ctx.db.Get())

	ctx.emailService = service.NewEmailService()
	ctx.subscriptionService = service.NewSubscriptionService(ctx.subscriptionRepository)

	ctx.rateStorage = service.NewRateStorage()

	ctx.rabbitMqConsumer = messaging.NewRabbitMQConsumer(getRabbitMQConnectionUrl())
	err := ctx.rabbitMqConsumer.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Printf("Connected to RabbitMQ")
	ctx.rateUpdateConsumer = consumers.NewRateUpdateConsumer(ctx.rateStorage)
	ctx.addSubscriptionConsumer = consumers.NewAddSubscriptionConsumer(ctx.subscriptionService)
	ctx.removeSubscriptionConsumer = consumers.NewRemoveSubscriptionConsumer(ctx.subscriptionService)
	ctx.consumerController = consumers.NewConsumerController(
		ctx.rabbitMqConsumer,
		ctx.rateUpdateConsumer,
		ctx.addSubscriptionConsumer,
		ctx.removeSubscriptionConsumer,
		getRabbitMqQueueName(),
	)

	ctx.consumerController.StartConsuming()
	ctx.subscriptionStatusPublisher = publishers.NewSubscriptionPublisher(
		ctx.rabbitMqConsumer,
		getRabbitMqQueueName(),
	)

	ctx.sendEmailJob = jobs.NewSendEmailJob(ctx.rateStorage, ctx.subscriptionService, ctx.emailService)

	ctx.initCron()
}

func (ctx *AppContext) Close() {
	ctx.db.Close()
	ctx.cronInstance.Stop()
	err := ctx.rabbitMqConsumer.Close()
	if err != nil {
		log.Printf("Failed to close RabbitMQ connection: %v", err)
	}
}

func getDbConfig() db.DbConfig {
	return db.DbConfig{
		Host:     util.GetEnv("DB_HOST", "localhost"),
		Port:     util.GetEnv("DB_PORT", "5432"),
		User:     util.GetEnv("DB_USER", "root"),
		Password: util.GetEnv("DB_PASSWORD", "password"),
		Name:     util.GetEnv("DB_NAME", "notification_service"),
		SSLMode:  util.GetEnv("DB_SSLMODE", "disable"),
	}
}

func (ctx *AppContext) initCron() {
	ctx.cronInstance = cron.New()

	_, err := ctx.cronInstance.AddFunc(util.GetEnv("EMAIL_CRON_CONFIG", "@daily"), ctx.sendEmailJob.SendEmails)
	if err != nil {
		log.Fatal()
	}

	ctx.cronInstance.Start()
}

func getRabbitMQConnectionUrl() string {
	return util.GetEnv("RABBITMQ_URL", "amqp://username:password@localhost:5672/")
}

func getRabbitMqQueueName() string {
	return util.GetEnv("RABBITMQ_QUEUE_NAME", "currency-notifier")
}
