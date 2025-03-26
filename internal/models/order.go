package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderID string `json:"order_id" gorm:"uniqueIndex"`
	Status  string `json:"status"`
}
