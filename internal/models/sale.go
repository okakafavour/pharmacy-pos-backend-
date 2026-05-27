package models

import "gorm.io/gorm"

type Sale struct {
	gorm.Model

	UserID uint    `json:"user_id"`
	Total  float64 `json:"total"`

	User      User       `gorm:"foreignKey:UserID"`
	SaleItems []SaleItem `gorm:"foreignKey:SaleID"`
}
