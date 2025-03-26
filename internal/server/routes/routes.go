package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/server/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(router *gin.Engine, orderHandler *handler.OrderHandler) {
	router.GET("/api/orders/:id", orderHandler.GetOrder)

	// Expose Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
