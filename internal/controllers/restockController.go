package controllers

import (
	"net/http"

	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type RestockInput struct {
	MedicineID uint `json:"medicine_id"`
	SupplierID uint `json:"supplier_id"`
	Quantity   int  `json:"quantity"`
}

func CreateRestock(c *gin.Context) {

	var input RestockInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var medicine models.Medicine

	config.DB.First(&medicine, input.MedicineID)

	if medicine.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Medicine not found",
		})
		return
	}

	var supplier models.Supplier

	config.DB.First(&supplier, input.SupplierID)

	if supplier.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Supplier not found",
		})
		return
	}

	restock := models.Restock{
		MedicineID: input.MedicineID,
		SupplierID: input.SupplierID,
		Quantity:   input.Quantity,
	}

	config.DB.Create(&restock)

	medicine.Stock += input.Quantity

	config.DB.Save(&medicine)

	config.DB.Preload("Medicine").
		Preload("Supplier").
		First(&restock, restock.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Medicine restocked successfully",
		"data":    restock,
	})
}
