package controllers

import (
	"net/http"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SaleItemInput struct {
	MedicineID uint `json:"medicine_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required"`
}

type CreateSaleInput struct {
	Items []SaleItemInput `json:"items" binding:"required,dive"`
}

func CreateSale(c *gin.Context) {
	var input CreateSaleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, _ := c.Get("user")
	claims := userData.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var total float64

	tx := config.DB.Begin()

	sale := models.Sale{
		UserID: userID,
		Total:  0,
	}

	tx.Create(&sale)

	for _, item := range input.Items {
		var medicine models.Medicine
		tx.First(&medicine, item.MedicineID)

		if medicine.ID == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Medicine not Found"})
			return
		}

		if medicine.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for " + medicine.Name})
			return
		}

		itemTotal := medicine.Price * float64(item.Quantity)
		total += itemTotal

		medicine.Stock -= item.Quantity
		tx.Save(&medicine)

		saleItem := models.SaleItem{
			SaleID:     sale.ID,
			MedicineID: medicine.ID,
			Quantity:   item.Quantity,
			Price:      medicine.Price,
		}

		tx.Create(&saleItem)
	}

	sale.Total = total
	tx.Save(&sale)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Sale created successfully",
		"sale_id": sale.ID,
	})
}
