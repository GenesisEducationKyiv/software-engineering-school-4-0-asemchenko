package main

import (
	"go.uber.org/zap"
	"notification-service/internal/context"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	appCtx := context.NewAppContext()
	appCtx.Init()

	defer appCtx.Close()
}
