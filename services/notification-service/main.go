package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"notification-service/internal/context"
	"notification-service/internal/metrics"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		zap.S().Fatalf("Failed to create logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			zap.S().Errorf("Failed to sync logger: %v", err)
		}
	}(logger)
	zap.ReplaceGlobals(logger)

	metrics.Init()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		zap.S().Info("Starting metrics server at :8082")
		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			zap.S().Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	appCtx := context.NewAppContext()
	appCtx.Init()

	defer appCtx.Close()
}
