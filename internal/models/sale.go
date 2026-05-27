package models

import "gorm.io/gorm"

type Sale struct {
	gorm.Model

	UserID     uint    `json:"user_id"`
	CustomerID uint    `json:"customer_id"`
	Total      float64 `json:"total"`

	User      User       `json:"user"`
	Customer  Customer   `json:"customer"`
	SaleItems []SaleItem `json:"sale_items"`
}
