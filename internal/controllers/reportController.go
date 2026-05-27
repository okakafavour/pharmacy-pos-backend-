package controllers

import (
	"net/http"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDailySalesReport(c *gin.Context) {

	var sales []models.Sale

	now := time.Now()

	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	endOfDay := startOfDay.Add(24 * time.Hour)

	config.DB.
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Find(&sales)

	var totalSales float64

	for _, sale := range sales {
		totalSales += sale.Total
	}

	c.JSON(http.StatusOK, gin.H{
		"date":               now.Format("2006-01-02"),
		"total_sales":        totalSales,
		"total_transactions": len(sales),
	})
}

func GetTopMedicines(c *gin.Context) {

	type Result struct {
		MedicineName string
		TotalSold    int
	}

	var results []Result

	config.DB.Raw(`
		SELECT medicines.name AS medicine_name,
		SUM(sale_items.quantity) AS total_sold
		FROM sale_items
		JOIN medicines
		ON medicines.id = sale_items.medicine_id
		GROUP BY medicines.name
		ORDER BY total_sold DESC
		LIMIT 5
	`).Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
