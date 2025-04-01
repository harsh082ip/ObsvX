// internal/repositories/order_repository.go
package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/harsh082ip/ObsvX/internal/metrics"
	"github.com/harsh082ip/ObsvX/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) FindByID(orderID string) (*models.Order, error) {
	start := time.Now()

	var order models.Order
	err := r.DB.Where("order_id = ?", orderID).First(&order).Error

	// Record database query duration
	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.WithLabelValues("find", "select").Observe(duration)
	metrics.DatabaseQueriesPerRequest.WithLabelValues("/api/orders/:id").Inc()

	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Create(req *models.CreateOrderRequest) (*models.Order, error) {
	start := time.Now()

	order := &models.Order{
		OrderID:     uuid.New().String(),
		Status:      "pending",
		Description: req.Description,
		Amount:      req.Amount,
	}

	err := r.DB.Create(order).Error

	// Record database query duration
	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.WithLabelValues("create", "insert").Observe(duration)
	metrics.DatabaseQueriesPerRequest.WithLabelValues("/api/orders").Inc()

	if err != nil {
		return nil, err
	}
	return order, nil
}
