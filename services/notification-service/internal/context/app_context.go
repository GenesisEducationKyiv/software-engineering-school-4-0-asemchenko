package context

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-asemchenko/shared/util"
	"github.com/robfig/cron/v3"
)

type AppContext struct {
	db           *db.Db
	cronInstance *cron.Cron
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

func (ctx *AppContext) Init() {
	ctx.db = db.NewDb(getDbConfig())
	ctx.db.Init()
}

func getDbConfig() db.DbConfig {
	return db.DbConfig{
		Host:     util.GetEnv("DB_HOST", "localhost"),
		Port:     util.GetEnv("DB_PORT", "5432"),
		User:     util.GetEnv("DB_USER", "postgres"),
		Password: util.GetEnv("DB_PASSWORD", "postgres"),
		Name:     util.GetEnv("DB_NAME", "notification_service"),
		SSLMode:  util.GetEnv("DB_SSLMODE", "disable"),
	}
}

// TODO anse
//_, err = ctx.cronInstance.AddFunc("@daily", ctx.SendEmailJob.SendEmails)
//if err != nil {
//log.Fatal()
//}
