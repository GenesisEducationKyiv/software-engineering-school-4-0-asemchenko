package main

import (
	"currency-notifier/internal/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Маршруты для контроллеров
	r.HandleFunc("/api/rate", controller.GetRate).Methods("GET")
	r.HandleFunc("/api/subscribe", controller.Subscribe).Methods("POST")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
