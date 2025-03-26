package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/config"
	"github.com/harsh082ip/ObsvX/internal/db"
	"github.com/harsh082ip/ObsvX/internal/metrics"
	"github.com/harsh082ip/ObsvX/internal/repositories"
	"github.com/harsh082ip/ObsvX/internal/server/handler"
	"github.com/harsh082ip/ObsvX/internal/server/routes"
	"github.com/harsh082ip/ObsvX/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Database
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	metrics.InitMetrics()
	orderRepo := repositories.NewOrderRepository(dbConn)
	orderService := services.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	router := gin.Default()
	routes.SetupRoutes(router, orderHandler)

	log.Println("🚀 Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
