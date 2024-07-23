package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"notification-service/internal/context"
	"notification-service/internal/metrics"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
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
