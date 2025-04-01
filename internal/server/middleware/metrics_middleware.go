package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/metrics"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}

		start := time.Now()

		// Increment the concurrent requests before processing
		metrics.ConcurrentRequests.WithLabelValues(endpoint).Inc()

		// Process request
		c.Next()

		// Decrement the concurrent requests after processing
		metrics.ConcurrentRequests.WithLabelValues(endpoint).Dec()

		// Record latency
		latency := time.Since(start).Seconds()
		metrics.ApiLatency.WithLabelValues(endpoint).Observe(latency)

		// Record response code
		statusCode := strconv.Itoa(c.Writer.Status())
		metrics.HttpResponseCodes.WithLabelValues(endpoint, statusCode, c.Request.Method).Inc()

	}
}
