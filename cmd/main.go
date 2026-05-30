package main

import (
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/controllers"
	"pharmacy-pos-backend/internal/models"
	"pharmacy-pos-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Medicine{},
		&models.Sale{},
		&models.SaleItem{},
		&models.Supplier{},
		&models.Restock{},
		&models.Customer{},
	)
	controllers.SeedAdmin()

	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":8080")
}
