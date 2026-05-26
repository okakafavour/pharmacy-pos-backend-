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

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST(
		"/medicines",
		middleware.AuthMiddleware(),
		middleware.RoleMiddleware("admin"),
		controllers.CreateMedicine,
	)
	router.GET("/medicines", controllers.GetMedicines)
}
