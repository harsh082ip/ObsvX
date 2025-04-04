package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/log"
	"github.com/harsh082ip/ObsvX/internal/metrics"
	"github.com/harsh082ip/ObsvX/internal/models"
	"github.com/harsh082ip/ObsvX/internal/services"
)

type OrderHandler struct {
	Service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	logger := log.InitLogger("handler")
	logger.DefaultLogger.RemoteAddr = c.ClientIP()

	start := time.Now()
	orderID := c.Param("id")
	endpoint := "/api/orders/:id"

	logger.LogInfoMessage().
		Str("operation", "GetOrder").
		Str("order_id", orderID).
		Str("endpoint", endpoint).
		Str("method", "GET").
		Msg("Processing get order request")

	metrics.RequestCount.WithLabelValues(endpoint, "GET").Inc()
	metrics.RequestsPerSecond.WithLabelValues(endpoint, "GET").Inc()

	order, err := h.Service.GetOrderByID(orderID)
	if err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "404").Inc()
		metrics.HttpResponseCodes.WithLabelValues(endpoint, "404", "GET").Inc()

		logger.LogWarnMessage().
			Str("error", err.Error()).
			Str("operation", "GetOrder").
			Str("order_id", orderID).
			Str("endpoint", endpoint).
			Int("status", http.StatusNotFound).
			Dur("latency", time.Since(start)).
			Msg("Order not found")

		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	latency := time.Since(start)
	metrics.ApiLatency.WithLabelValues(endpoint).Observe(latency.Seconds())
	metrics.HttpResponseCodes.WithLabelValues(endpoint, "200", "GET").Inc()

	logger.LogInfoMessage().
		Str("operation", "GetOrder").
		Str("order_id", orderID).
		Str("endpoint", endpoint).
		Int("status", http.StatusOK).
		Dur("latency", latency).
		Msg("Order retrieved successfully")

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	logger := log.InitLogger("handler")
	logger.DefaultLogger.RemoteAddr = c.ClientIP()

	start := time.Now()
	endpoint := "/api/orders"

	logger.LogInfoMessage().
		Str("operation", "CreateOrder").
		Str("endpoint", endpoint).
		Str("method", "POST").
		Msg("Processing create order request")

	metrics.RequestCount.WithLabelValues(endpoint, "POST").Inc()
	metrics.RequestsPerSecond.WithLabelValues(endpoint, "POST").Inc()

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "400").Inc()
		metrics.HttpResponseCodes.WithLabelValues(endpoint, "400", "POST").Inc()

		logger.LogWarnMessage().
			Str("error", err.Error()).
			Str("operation", "CreateOrder").
			Str("endpoint", endpoint).
			Int("status", http.StatusBadRequest).
			Dur("latency", time.Since(start)).
			Msg("Invalid request payload")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Service.CreateOrder(&req)
	if err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "500").Inc()
		metrics.HttpResponseCodes.WithLabelValues(endpoint, "500", "POST").Inc()

		logger.LogErrorMessage().
			Str("error", err.Error()).
			Str("operation", "CreateOrder").
			Str("endpoint", endpoint).
			Float64("amount", req.Amount).
			Str("description", req.Description).
			Int("status", http.StatusInternalServerError).
			Dur("latency", time.Since(start)).
			Msg("Failed to create order")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	latency := time.Since(start)
	metrics.ApiLatency.WithLabelValues(endpoint).Observe(latency.Seconds())
	metrics.HttpResponseCodes.WithLabelValues(endpoint, "201", "POST").Inc()

	// Extract the order ID for logging if possible
	orderID := "unknown"
	if orderObj, ok := order.(*models.Order); ok {
		orderID = orderObj.OrderID
	}

	logger.LogInfoMessage().
		Str("operation", "CreateOrder").
		Str("order_id", orderID).
		Str("endpoint", endpoint).
		Float64("amount", req.Amount).
		Str("description", req.Description).
		Int("status", http.StatusCreated).
		Dur("latency", latency).
		Msg("Order created successfully")

	c.JSON(http.StatusCreated, order)
}
