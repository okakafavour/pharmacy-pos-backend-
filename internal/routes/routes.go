package routes

import (
	"pharmacy-pos-backend/internal/controllers"
	middleware "pharmacy-pos-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pharmacy POS API Running",
		})
	})

	// Public Routes
	router.POST("/login", controllers.Login)

	// Protected Routes
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Any logged in user can view medicines
		authorized.GET("/medicines", controllers.GetMedicines)
		authorized.POST("/sales", controllers.CreateSale)

		// Admin Routes
		admin := authorized.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("/users", controllers.CreateUser)

			admin.POST("/medicines", controllers.CreateMedicine)
		}
	}
}
