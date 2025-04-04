package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.InitLogger("http")
		logger.DefaultLogger.RemoteAddr = c.ClientIP()

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate request duration
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details
		logger.LogInfoMessage().
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Msg("Request processed")
	}
}
