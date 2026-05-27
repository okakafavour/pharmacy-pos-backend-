package models

import "gorm.io/gorm"

type SaleItem struct {
	gorm.Model

	SaleID     uint    `json:"sale_id"`
	MedicineID uint    `json:"medicine_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`

	Medicine Medicine `gorm:"foreignKey:MedicineID"`
}
