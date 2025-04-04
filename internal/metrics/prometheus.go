package metrics

import (
	"net/http"
	"runtime"
	"time"

	"github.com/harsh082ip/ObsvX/internal/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create a custom registry
var Registry = prometheus.NewRegistry()

var (
	RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_request_count",
		Help: "Total API requests",
	}, []string{"endpoint", "method"})

	ApiLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_latency_seconds",
		Help:    "API response time",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"endpoint"})

	ApiErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_errors_total",
		Help: "Total API errors",
	}, []string{"endpoint", "status_code"})

	ConcurrentRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "api_concurrent_requests",
		Help: "Number of concurrent requests being processed",
	}, []string{"endpoint"})

	RequestsPerSecond = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_requests_per_second",
		Help: "Number of requests per second",
	}, []string{"endpoint", "method"})

	DatabaseQueryDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_query_duration_seconds",
		Help:    "Duration of database queries in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	}, []string{"query_type", "operation"})

	DatabaseQueriesPerRequest = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "database_queries_per_request",
		Help: "Number of database queries per request",
	}, []string{"endpoint"})

	HttpResponseCodes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_response_codes_total",
		Help: "Count of HTTP response codes",
	}, []string{"endpoint", "code", "method"})

	// CPU and Memory metrics
	CpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_cpu_usage_percent",
		Help: "Application CPU usage in percent",
	})

	MemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_memory_usage_bytes",
		Help: "Application memory usage in bytes",
	})
)

// InitMetrics registers metrics with the custom registry
func InitMetrics() {
	logger := log.InitLogger("metrics")

	Registry.MustRegister(
		RequestCount,
		ApiLatency,
		ApiErrors,
		ConcurrentRequests,
		RequestsPerSecond,
		DatabaseQueryDuration,
		DatabaseQueriesPerRequest,
		HttpResponseCodes,
		CpuUsage,
		MemoryUsage,
	)

	logger.LogInfoMessage().Msg("Metrics registered successfully")

	// Start a goroutine to periodically update CPU and memory metrics
	go updateResourceMetrics()
}

func updateResourceMetrics() {
	logger := log.InitLogger("metrics")
	var memStats runtime.MemStats

	logger.LogDebugMessage().Msg("Starting resource metrics collection")

	for {
		runtime.ReadMemStats(&memStats)
		MemoryUsage.Set(float64(memStats.Alloc))

		CpuUsage.Set(float64(runtime.NumGoroutine()) / 100)

		// Update every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func Handler() http.Handler {
	return promhttp.HandlerFor(Registry, promhttp.HandlerOpts{})
}
