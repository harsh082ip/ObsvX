package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/metrics"
	"github.com/harsh082ip/ObsvX/internal/models"
	"github.com/harsh082ip/ObsvX/internal/services"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	Service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	start := time.Now()
	orderID := c.Param("id")
	endpoint := "/api/orders/:id"

	metrics.RequestCount.WithLabelValues(endpoint, "GET").Inc()
	metrics.RequestsPerSecond.WithLabelValues(endpoint, "GET").Inc()

	order, err := h.Service.GetOrderByID(orderID)
	if err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "404").Inc()
		metrics.HttpResponseCodes.WithLabelValues(endpoint, "404", "GET").Inc()
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	latency := time.Since(start).Seconds()
	metrics.ApiLatency.WithLabelValues(endpoint).Observe(latency)
	metrics.HttpResponseCodes.WithLabelValues(endpoint, "200", "GET").Inc()

	logrus.WithFields(logrus.Fields{"order_id": orderID, "latency": latency}).Info("Request successful")

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	start := time.Now()
	endpoint := "/api/orders"

	metrics.RequestCount.WithLabelValues(endpoint, "POST").Inc()
	metrics.RequestsPerSecond.WithLabelValues(endpoint, "POST").Inc()

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "400").Inc()
		metrics.HttpResponseCodes.WithLabelValues(endpoint, "400", "POST").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Service.CreateOrder(&req)
	if err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "500").Inc()
		metrics.HttpResponseCodes.WithLabelValues(endpoint, "500", "POST").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	latency := time.Since(start).Seconds()
	metrics.ApiLatency.WithLabelValues(endpoint).Observe(latency)
	metrics.HttpResponseCodes.WithLabelValues(endpoint, "201", "POST").Inc()

	logrus.WithFields(logrus.Fields{
		"latency": latency,
		"order":   order,
	}).Info("Order created successfully")

	c.JSON(http.StatusCreated, order)
}
