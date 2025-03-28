package services

import (
	"errors"

	"github.com/harsh082ip/ObsvX/internal/repositories"
)

type OrderService struct {
	Repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{Repo: repo}
}

func (s *OrderService) GetOrderByID(orderID string) (interface{}, error) {
	order, err := s.Repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}
