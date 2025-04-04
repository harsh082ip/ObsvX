package services

import (
	"errors"
	"time"

	"github.com/harsh082ip/ObsvX/internal/log"
	"github.com/harsh082ip/ObsvX/internal/models"
	"github.com/harsh082ip/ObsvX/internal/repositories"
)

type OrderService struct {
	Repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{Repo: repo}
}

func (s *OrderService) GetOrderByID(orderID string) (interface{}, error) {
	logger := log.InitLogger("service")
	start := time.Now()

	logger.LogDebugMessage().
		Str("operation", "GetOrderByID").
		Str("order_id", orderID).
		Msg("Getting order by ID")

	order, err := s.Repo.FindByID(orderID)

	if err != nil {
		logger.LogDebugMessage().
			Str("error", err.Error()).
			Str("operation", "GetOrderByID").
			Str("order_id", orderID).
			Dur("duration", time.Since(start)).
			Msg("Order not found")
		return nil, errors.New("order not found")
	}

	logger.LogDebugMessage().
		Str("operation", "GetOrderByID").
		Str("order_id", orderID).
		Dur("duration", time.Since(start)).
		Msg("Order retrieved successfully")

	return order, nil
}

func (s *OrderService) CreateOrder(req *models.CreateOrderRequest) (interface{}, error) {
	logger := log.InitLogger("service")
	start := time.Now()

	logger.LogDebugMessage().
		Str("operation", "CreateOrder").
		Float64("amount", req.Amount).
		Str("description", req.Description).
		Msg("Creating order")

	order, err := s.Repo.Create(req)

	if err != nil {
		logger.LogErrorMessage().
			Str("error", err.Error()).
			Str("operation", "CreateOrder").
			Float64("amount", req.Amount).
			Str("description", req.Description).
			Dur("duration", time.Since(start)).
			Msg("Failed to create order")
		return nil, errors.New("failed to create order")
	}

	logger.LogDebugMessage().
		Str("operation", "CreateOrder").
		Str("order_id", order.OrderID).
		Float64("amount", req.Amount).
		Str("description", req.Description).
		Dur("duration", time.Since(start)).
		Msg("Order created successfully")

	return order, nil
}
