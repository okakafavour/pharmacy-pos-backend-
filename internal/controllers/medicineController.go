package controllers

import (
	"net/http"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateMedicine(c *gin.Context) {
	var medicine models.Medicine
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Create(&medicine)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Medicine created successfully",
		"data":    medicine,
	})
}

func GetMedicines(c *gin.Context) {
	var medicines []models.Medicine
	config.DB.Find(&medicines)
	c.JSON(http.StatusOK, gin.H{
		"data": medicines,
	})
}

func GetLowStockMedicines(c *gin.Context) {

	var medicines []models.Medicine

	config.DB.
		Where("stock < ?", 10).
		Find(&medicines)

	c.JSON(http.StatusOK, gin.H{
		"data": medicines,
	})
}

func GetExpiredMedicines(c *gin.Context) {

	var medicines []models.Medicine

	config.DB.
		Where("expire_date < NOW()").
		Find(&medicines)

	c.JSON(http.StatusOK, gin.H{
		"data": medicines,
	})
}

func GetExpiringSoonMedicines(c *gin.Context) {

	var medicines []models.Medicine

	config.DB.
		Where("expire_date BETWEEN NOW() AND NOW() + INTERVAL '30 days'").
		Find(&medicines)

	c.JSON(http.StatusOK, gin.H{
		"data": medicines,
	})
}
