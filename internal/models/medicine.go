package models

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	gorm.Model
	Name        string    `json:"name"`
	Barcode     string    `json:"barcode" gorm:"unique"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ExpireDate  time.Time `json:"expire_date"`
}
