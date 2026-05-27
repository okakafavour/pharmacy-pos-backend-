package controllers

import (
	"net/http"

	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {

	var totalMedicines int64
	var lowStockCount int64
	var expiredCount int64
	var totalSuppliers int64

	var totalSales float64

	config.DB.Model(&models.Medicine{}).Count(&totalMedicines)

	config.DB.Model(&models.Medicine{}).
		Where("stock < ?", 10).
		Count(&lowStockCount)

	config.DB.Model(&models.Medicine{}).
		Where("expire_date < NOW()").
		Count(&expiredCount)

	config.DB.Model(&models.Supplier{}).
		Count(&totalSuppliers)

	config.DB.Model(&models.Sale{}).
		Select("COALESCE(SUM(total),0)").
		Scan(&totalSales)

	c.JSON(http.StatusOK, gin.H{
		"total_medicines": totalMedicines,
		"low_stock_count": lowStockCount,
		"expired_count":   expiredCount,
		"total_sales":     totalSales,
		"total_suppliers": totalSuppliers,
	})
}
