package main

import (
	_ "currency-notifier/docs"
	"currency-notifier/internal/context"
	"currency-notifier/internal/server"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title UAH currency application
// @version 1.0
// @description API for current USD-UAH exchange rate and for email-subscribing on the currency rate
// @host localhost:8080
// @basePath /api

//go:generate go run github.com/swaggo/swag/cmd/swag init

func main() {
	appCtx := context.NewAppContext()

	appCtx.Init()

	s := server.NewServer(appCtx)
	s.RegisterRoutes()
	go s.StartListening()

	// graceful shutdown
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish
	log.Printf("Shutting down the server...")
	appCtx.Close()
}
