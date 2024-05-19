package main

import (
	"currency-notifier/internal/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "currency-notifier/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title UAH currency application
// @version 1.0
// @description API for current USD-UAH exchange rate and for email-subscribing on the currency rate
// @host localhost:8080
// @basePath /api

//go:generate go run github.com/swaggo/swag/cmd/swag init

func main() {
	r := mux.NewRouter()

	// Маршруты для контроллеров
	r.HandleFunc("/api/v1/rate", controller.GetRate).Methods("GET")
	r.HandleFunc("/api/v1/subscribe", controller.Subscribe).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
