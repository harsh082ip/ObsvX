package repositories

import (
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
	var order models.Order
	if err := r.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
