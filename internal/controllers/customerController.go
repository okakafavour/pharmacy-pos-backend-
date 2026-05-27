package controllers

import (
	"net/http"

	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CustomerInput struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func CreateCustomer(c *gin.Context) {

	var input CustomerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	customer := models.Customer{
		Name:    input.Name,
		Phone:   input.Phone,
		Address: input.Address,
	}

	config.DB.Create(&customer)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"data":    customer,
	})
}

func GetCustomers(c *gin.Context) {

	var customers []models.Customer

	config.DB.Find(&customers)

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
	})
}
