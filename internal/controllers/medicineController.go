package controllers

import (
	"net/http"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"
	"strconv"

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

func UpdateMedicine(c *gin.Context) {

	id := c.Param("id")

	var medicine models.Medicine

	err := config.DB.First(&medicine, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Medicine not found",
		})
		return
	}

	var input models.Medicine

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	medicine.Name = input.Name
	medicine.Barcode = input.Barcode
	medicine.Description = input.Description
	medicine.Price = input.Price
	medicine.Stock = input.Stock
	medicine.ExpireDate = input.ExpireDate

	config.DB.Save(&medicine)

	c.JSON(http.StatusOK, gin.H{
		"message": "Medicine updated successfully",
		"data":    medicine,
	})
}

func DeleteMedicine(c *gin.Context) {

	id := c.Param("id")

	var medicine models.Medicine

	err := config.DB.First(&medicine, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Medicine not found",
		})
		return
	}

	config.DB.Delete(&medicine)

	c.JSON(http.StatusOK, gin.H{
		"message": "Medicine deleted successfully",
	})
}

func GetMedicines(c *gin.Context) {

	var medicines []models.Medicine

	// Query Params
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Convert page and limit to int
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	offset := (pageInt - 1) * limitInt

	query := config.DB.Model(&models.Medicine{})

	// Search
	if search != "" {
		query = query.Where(
			"name ILIKE ? OR barcode ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	// Pagination
	query.
		Limit(limitInt).
		Offset(offset).
		Find(&medicines)

	c.JSON(http.StatusOK, gin.H{
		"page":  pageInt,
		"limit": limitInt,
		"data":  medicines,
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

func GetMedicineByBarcode(c *gin.Context) {

	barcode := c.Param("barcode")

	var medicine models.Medicine

	err := config.DB.
		Where("barcode = ?", barcode).
		First(&medicine).Error

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Medicine not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": medicine,
	})
}
