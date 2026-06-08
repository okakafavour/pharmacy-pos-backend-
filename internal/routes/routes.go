package routes

import (
	"pharmacy-pos-backend/internal/controllers"
	middleware "pharmacy-pos-backend/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://pharmacy-pos-frontend-bice.vercel.app",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pharmacy POS API Running",
		})
	})

	// Public Routes
	router.POST("/login", controllers.Login)
	router.POST("/google/login", controllers.GoogleLogin)

	// Protected Routes
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	{
		authorized.GET("/medicines", controllers.GetMedicines)
		authorized.GET("/medicines/barcode/:barcode", controllers.GetMedicineByBarcode)

		authorized.POST("/customers", controllers.CreateCustomer)
		authorized.GET("/customers", controllers.GetCustomers)

		authorized.GET("/sales", controllers.GetSales)
		authorized.GET("/sales/:id/receipt", controllers.GetReceipt)
		authorized.GET("/receipt/:id/pdf", controllers.DownloadReceiptPDF)

		sales := authorized.Group("/")
		sales.Use(middleware.RequireRole("admin", "cashier"))
		{
			sales.POST("/sales", controllers.CreateSale)
		}

		medicine := authorized.Group("/")
		medicine.Use(middleware.RequireRole("admin", "pharmacist"))
		{
			medicine.POST("/medicines", controllers.CreateMedicine)
			medicine.PUT("/medicines/:id", controllers.UpdateMedicine)
			medicine.DELETE("/medicines/:id", controllers.DeleteMedicine)

			medicine.GET("/medicines/low-stock", controllers.GetLowStockMedicines)
			medicine.GET("/medicines/expired", controllers.GetExpiredMedicines)
			medicine.GET("/medicines/expiring-soon", controllers.GetExpiringSoonMedicines)
		}

		reports := authorized.Group("/")
		reports.Use(middleware.RequireRole("admin", "manager"))
		{
			reports.GET("/dashboard/stats", controllers.GetDashboardStats)
			reports.GET("/reports/daily-sales", controllers.GetDailySalesReport)
			reports.GET("/reports/top-medicines", controllers.GetTopMedicines)
		}

		suppliers := authorized.Group("/")
		suppliers.Use(middleware.RequireRole("admin", "manager"))
		{
			suppliers.POST("/suppliers", controllers.CreateSupplier)
			suppliers.GET("/suppliers", controllers.GetSuppliers)
			suppliers.POST("/restocks", controllers.CreateRestock)
			suppliers.GET("/restocks", controllers.GetRestocks)
		}

		admin := authorized.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("/users", controllers.CreateUser)
		}
	}
}
