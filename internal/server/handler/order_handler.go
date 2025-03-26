package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/ObsvX/internal/metrics"
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

	order, err := h.Service.GetOrderByID(orderID)
	if err != nil {
		metrics.ApiErrors.WithLabelValues(endpoint, "404").Inc()
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	latency := time.Since(start).Seconds()
	metrics.ApiLatency.WithLabelValues(endpoint).Observe(latency)

	logrus.WithFields(logrus.Fields{"order_id": orderID, "latency": latency}).Info("Request successful")

	c.JSON(http.StatusOK, order)
}
