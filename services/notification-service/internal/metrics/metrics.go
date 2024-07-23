package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	NotificationsSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_sent_total",
			Help: "Total number of sent notifications.",
		},
		[]string{"type"},
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
		NotificationsSent,
		RabbitMQSentMessages,
		RabbitMQReceivedMessages,
	)
}
