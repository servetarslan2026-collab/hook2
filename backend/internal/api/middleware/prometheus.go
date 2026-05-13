package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
		},
		[]string{"method", "path"},
	)

	webhookDeliveriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webhook_deliveries_total",
			Help: "Total webhook delivery attempts",
		},
		[]string{"status"},
	)

	activeWorkers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "webhook_active_workers",
			Help: "Number of active webhook delivery workers",
		},
	)

	eventsReceived = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "events_received_total",
			Help: "Total events received via API",
		},
	)
)

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		path := c.Route().Path
		method := c.Method()
		status := strconv.Itoa(c.Response().StatusCode())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

func IncWebhookDelivery(status string) {
	webhookDeliveriesTotal.WithLabelValues(status).Inc()
}

func IncEventsReceived() {
	eventsReceived.Inc()
}

func SetActiveWorkers(n float64) {
	activeWorkers.Set(n)
}
