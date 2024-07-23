package server

import (
	"currency-notifier/internal/context"
	"currency-notifier/internal/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
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

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) RegisterRoutes() {
	subscriptionController := controller.NewSubscriptionController(s.ctx.UserService)
	rateController := controller.NewRateController(s.ctx.CurrencyService)

	s.router.HandleFunc("/api/rate", rateController.GetRate).Methods("GET")
	s.router.HandleFunc("/api/subscribe", subscriptionController.Subscribe).Methods("POST")
	s.router.HandleFunc("/api/unsubscribe", subscriptionController.Unsubscribe).Methods("POST")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

func (s *Server) StartListening() {
	zap.S().Info("Server started at :8080")
	zap.S().Fatal(http.ListenAndServe(":8080", s.router))
}
