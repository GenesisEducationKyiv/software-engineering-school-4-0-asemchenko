package main

import (
	_ "currency-notifier/docs"
	"currency-notifier/internal"
	"currency-notifier/internal/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title UAH currency application
// @version 1.0
// @description API for current USD-UAH exchange rate and for email-subscribing on the currency rate
// @host localhost:8080
// @basePath /api

//go:generate go run github.com/swaggo/swag/cmd/swag init

func main() {
	r := mux.NewRouter()
	appCtx := internal.NewAppContext()

	appCtx.Init()
	defer appCtx.Close()

	subscriptionController := controller.NewSubscriptionController(appCtx.SubscriptionService)
	rateController := controller.NewRateController(appCtx.CurrencyService)

	r.HandleFunc("/api/rate", rateController.GetRate).Methods("GET")
	r.HandleFunc("/api/subscribe", subscriptionController.Subscribe).Methods("POST")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
