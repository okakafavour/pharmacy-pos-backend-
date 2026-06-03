package controllers

import (
	"log"
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

	log.Println("===== CREATE CUSTOMER START =====")

	var input CustomerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("JSON ERROR:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("INPUT:", input)

	customer := models.Customer{
		Name:    input.Name,
		Phone:   input.Phone,
		Address: input.Address,
	}

	log.Println("BEFORE DB CREATE")

	result := config.DB.Create(&customer)

	if result.Error != nil {
		log.Println("DB ERROR:", result.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	log.Println("CUSTOMER CREATED:", customer.ID)

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
