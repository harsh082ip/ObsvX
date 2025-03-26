package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_request_count",
		Help: "Total API requests",
	}, []string{"endpoint", "method"})

	ApiLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_latency_seconds",
		Help:    "API response time",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint"})

	ApiErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_errors_total",
		Help: "Total API errors",
	}, []string{"endpoint", "status_code"})
)

// InitMetrics registers metrics with Prometheus
func InitMetrics() {
	prometheus.MustRegister(RequestCount, ApiLatency, ApiErrors)
}
