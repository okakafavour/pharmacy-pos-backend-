package main

import (
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/controllers"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Medicine{},
	)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Pharmacy POS Backend API",
		})
	})
	router.POST("/register", controllers.Register)

	router.Run(":8080")
}
