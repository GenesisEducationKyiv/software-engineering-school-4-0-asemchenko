package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ApiRateRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_rate_requests_total",
			Help: "Total number of /api/rate requests.",
		},
		[]string{"method"},
	)

	ApiRateDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_rate_duration_seconds",
			Help: "Duration of /api/rate requests.",
		},
		[]string{"method"},
	)
	ApiSubscribeRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_subscribe_requests_total",
			Help: "Total number of /api/subscribe requests.",
		},
		[]string{"method"},
	)
	ApiSubscribeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_subscribe_duration_seconds",
			Help: "Duration of /api/subscribe requests.",
		},
		[]string{"method"},
	)
	ApiUnsubscribeRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_unsubscribe_requests_total",
			Help: "Total number of /api/unsubscribe requests.",
		},
		[]string{"method"},
	)
	ApiUnsubscribeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_unsubscribe_duration_seconds",
			Help: "Duration of /api/unsubscribe requests.",
		},
		[]string{"method"},
	)
	RabbitMQSentMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_sent_messages_total",
			Help: "Total number of sent RabbitMQ messages.",
		},
		[]string{"type"},
	)

	RabbitMQReceivedMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_received_messages_total",
			Help: "Total number of received RabbitMQ messages.",
		},
		[]string{"type"},
	)
)

func Init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(
		ApiRateRequestsTotal,
		ApiRateDuration,
		ApiSubscribeRequestsTotal,
		ApiSubscribeDuration,
		ApiUnsubscribeRequestsTotal,
		ApiUnsubscribeDuration,
		RabbitMQSentMessages,
		RabbitMQReceivedMessages,
	)
}
