package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model

	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
