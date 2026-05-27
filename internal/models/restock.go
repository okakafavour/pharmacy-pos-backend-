package models

import "gorm.io/gorm"

type Restock struct {
	gorm.Model

	MedicineID uint `json:"medicine_id"`
	SupplierID uint `json:"supplier_id"`
	Quantity   int  `json:"quantity"`

	Medicine Medicine `json:"medicine" gorm:"foreignKey:MedicineID;references:ID"`

	Supplier Supplier `json:"supplier" gorm:"foreignKey:SupplierID;references:ID"`
}
