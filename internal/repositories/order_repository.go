// internal/repositories/order_repository.go
package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/harsh082ip/ObsvX/internal/log"
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
	logger := log.InitLogger("repository")
	start := time.Now()

	logger.LogDebugMessage().
		Str("operation", "FindByID").
		Str("order_id", orderID).
		Msg("Finding order by ID")

	var order models.Order
	err := r.DB.Where("order_id = ?", orderID).First(&order).Error

	// Record database query duration
	duration := time.Since(start)
	metrics.DatabaseQueryDuration.WithLabelValues("find", "select").Observe(duration.Seconds())
	metrics.DatabaseQueriesPerRequest.WithLabelValues("/api/orders/:id").Inc()

	if err != nil {
		logger.LogErrorMessage().
			Str("error", err.Error()).
			Str("operation", "FindByID").
			Str("order_id", orderID).
			Dur("duration", duration).
			Msg("Failed to find order")
		return nil, err
	}

	logger.LogDebugMessage().
		Str("operation", "FindByID").
		Str("order_id", orderID).
		Dur("duration", duration).
		Msg("Order found successfully")

	return &order, nil
}

func (r *OrderRepository) Create(req *models.CreateOrderRequest) (*models.Order, error) {
	logger := log.InitLogger("repository")
	start := time.Now()

	orderID := uuid.New().String()

	logger.LogDebugMessage().
		Str("operation", "Create").
		Str("order_id", orderID).
		Float64("amount", req.Amount).
		Str("description", req.Description).
		Msg("Creating new order")

	order := &models.Order{
		OrderID:     orderID,
		Status:      "pending",
		Description: req.Description,
		Amount:      req.Amount,
	}

	err := r.DB.Create(order).Error

	// Record database query duration
	duration := time.Since(start)
	metrics.DatabaseQueryDuration.WithLabelValues("create", "insert").Observe(duration.Seconds())
	metrics.DatabaseQueriesPerRequest.WithLabelValues("/api/orders").Inc()

	if err != nil {
		logger.LogErrorMessage().
			Str("error", err.Error()).
			Str("operation", "Create").
			Str("order_id", orderID).
			Dur("duration", duration).
			Msg("Failed to create order")
		return nil, err
	}

	logger.LogDebugMessage().
		Str("operation", "Create").
		Str("order_id", orderID).
		Dur("duration", duration).
		Msg("Order created successfully")

	return order, nil
}
