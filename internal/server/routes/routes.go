package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/metrics"
	"github.com/harsh082ip/ObsvX/internal/server/handler"
	"github.com/harsh082ip/ObsvX/internal/server/middleware"
)

func SetupRoutes(router *gin.Engine, orderHandler *handler.OrderHandler) {
	// Apply metrics middleware to all routes
	router.Use(middleware.MetricsMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/api/orders/:id", orderHandler.GetOrder)
	router.POST("/api/orders", orderHandler.CreateOrder)

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(metrics.Handler()))
}
