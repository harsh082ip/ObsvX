package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderID     string  `json:"order_id" gorm:"uniqueIndex"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type CreateOrderRequest struct {
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}
