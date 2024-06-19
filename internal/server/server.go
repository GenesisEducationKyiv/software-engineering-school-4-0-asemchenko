package server

import (
	"currency-notifier/internal/context"
	"currency-notifier/internal/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

type Server struct {
	router *mux.Router
	ctx    *context.AppContext
}

func NewServer(ctx *context.AppContext) *Server {
	return &Server{
		router: mux.NewRouter(),
		ctx:    ctx,
	}
}

func (s *Server) RegisterRoutes() {
	subscriptionController := controller.NewSubscriptionController(s.ctx.SubscriptionService)
	rateController := controller.NewRateController(s.ctx.CurrencyService)

	s.router.HandleFunc("/api/rate", rateController.GetRate).Methods("GET")
	s.router.HandleFunc("/api/subscribe", subscriptionController.Subscribe).Methods("POST")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

func (s *Server) StartListening() {
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", s.router))
}
