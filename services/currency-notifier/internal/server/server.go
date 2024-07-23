package server

import (
	"currency-notifier/internal/context"
	"currency-notifier/internal/controller"
	"currency-notifier/internal/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	s.router.Handle("/api/rate", promhttp.InstrumentHandlerDuration(
		metrics.ApiRateDuration,
		promhttp.InstrumentHandlerCounter(metrics.ApiRateRequestsTotal, http.HandlerFunc(rateController.GetRate)),
	)).Methods("GET")

	s.router.Handle("/api/subscribe", promhttp.InstrumentHandlerDuration(
		metrics.ApiSubscribeDuration,
		promhttp.InstrumentHandlerCounter(metrics.ApiSubscribeRequestsTotal, http.HandlerFunc(subscriptionController.Subscribe)),
	)).Methods("POST")

	s.router.Handle("/api/unsubscribe", promhttp.InstrumentHandlerDuration(
		metrics.ApiUnsubscribeDuration,
		promhttp.InstrumentHandlerCounter(metrics.ApiUnsubscribeRequestsTotal, http.HandlerFunc(subscriptionController.Unsubscribe)),
	)).Methods("POST")

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

func (s *Server) RegisterMetrics() {
	s.router.Handle("/metrics", promhttp.Handler())
}

func (s *Server) StartListening() {
	zap.S().Info("Server started at :8080")
	zap.S().Fatal(http.ListenAndServe(":8080", s.router))
}
