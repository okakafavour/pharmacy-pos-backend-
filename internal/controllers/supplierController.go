package controllers

import (
	"net/http"

	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateSupplier(c *gin.Context) {

	var supplier models.Supplier

	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	config.DB.Create(&supplier)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Supplier created successfully",
		"data":    supplier,
	})
}

func GetSuppliers(c *gin.Context) {

	var suppliers []models.Supplier

	config.DB.Find(&suppliers)

	c.JSON(http.StatusOK, gin.H{
		"data": suppliers,
	})
}
