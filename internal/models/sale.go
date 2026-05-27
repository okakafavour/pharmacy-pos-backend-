package models

import "gorm.io/gorm"

type Sale struct {
	gorm.Model

	UserID uint    `json:"user_id"`
	Total  float64 `json:"total"`

	User      User       `json:"user" gorm:"foreignKey:UserID"`
	SaleItems []SaleItem `json:"sale_items" gorm:"foreignKey:SaleID"`
}
