package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/config"
	"github.com/harsh082ip/ObsvX/internal/db"
	"github.com/harsh082ip/ObsvX/internal/log"
	"github.com/harsh082ip/ObsvX/internal/metrics"
	"github.com/harsh082ip/ObsvX/internal/repositories"
	"github.com/harsh082ip/ObsvX/internal/server/handler"
	"github.com/harsh082ip/ObsvX/internal/server/middleware"
	"github.com/harsh082ip/ObsvX/internal/server/routes"
	"github.com/harsh082ip/ObsvX/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("failed to load config, err: ", err.Error())
		os.Exit(1)
	}

	// Initialize global logger
	if err := log.SetupGlobalLogger(cfg); err != nil {
		fmt.Println("failed to setup global logger, err: ", err.Error())
		os.Exit(1)
	}

	logger := log.InitLogger("main")

	logger.LogInfoMessage().Msg("Starting application")

	gin.SetMode(gin.ReleaseMode)

	// Initialize Database
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		logger.LogFatalMessage().
			Str("error", err.Error()).
			Msg("Database connection failed")
		os.Exit(1)
	}
	logger.LogInfoMessage().Msg("Database connection established")

	// Initialize metrics
	metrics.InitMetrics()
	logger.LogInfoMessage().Msg("Metrics initialized")

	// Initialize repositories, services, and handlers
	orderRepo := repositories.NewOrderRepository(dbConn)
	orderService := services.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	// Initialize router with logging middleware
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.MetricsMiddleware())

	// Setup routes
	routes.SetupRoutes(router, orderHandler)

	// Start server
	serverAddr := ":8080"
	logger.LogInfoMessage().
		Str("address", serverAddr).
		Msg("Server starting")

	if err := router.Run(serverAddr); err != nil {
		logger.LogFatalMessage().
			Str("error", err.Error()).
			Msg("Server failed to start")
		os.Exit(1)
	}
}
