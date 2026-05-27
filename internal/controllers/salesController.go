package controllers

import (
	"net/http"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jung-kurt/gofpdf"
)

type SaleItemInput struct {
	MedicineID uint `json:"medicine_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required"`
}

type CreateSaleInput struct {
	CustomerID uint `json:"customer_id"`

	Items []SaleItemInput `json:"items" binding:"required,dive"`
}

func CreateSale(c *gin.Context) {

	var input CreateSaleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userData, _ := c.Get("user")
	claims := userData.(jwt.MapClaims)

	userID := uint(claims["user_id"].(float64))

	var total float64

	tx := config.DB.Begin()

	sale := models.Sale{
		UserID:     userID,
		CustomerID: input.CustomerID,
		Total:      0,
	}

	tx.Create(&sale)

	for _, item := range input.Items {

		var medicine models.Medicine

		tx.First(&medicine, item.MedicineID)

		if medicine.ID == 0 {
			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Medicine not found",
			})
			return
		}

		if medicine.Stock < item.Quantity {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock for " + medicine.Name,
			})
			return
		}

		itemTotal := medicine.Price * float64(item.Quantity)

		total += itemTotal

		medicine.Stock -= item.Quantity

		tx.Save(&medicine)

		saleItem := models.SaleItem{
			SaleID:     sale.ID,
			MedicineID: medicine.ID,
			Quantity:   item.Quantity,
			Price:      medicine.Price,
		}

		tx.Create(&saleItem)
	}

	sale.Total = total

	tx.Save(&sale)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Sale created successfully",
		"sale_id": sale.ID,
	})
}

func GetSales(c *gin.Context) {

	var sales []models.Sale

	config.DB.
		Preload("User").
		Preload("Customer").
		Preload("SaleItems").
		Preload("SaleItems.Medicine").
		Find(&sales)

	c.JSON(http.StatusOK, gin.H{
		"data": sales,
	})
}

func GetReceipt(c *gin.Context) {

	id := c.Param("id")

	var sale models.Sale

	err := config.DB.
		Preload("User").
		Preload("Customer").
		Preload("SaleItems").
		Preload("SaleItems.Medicine").
		First(&sale, id).Error

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Sale not found",
		})
		return
	}

	var items []gin.H

	for _, item := range sale.SaleItems {

		items = append(items, gin.H{
			"medicine": item.Medicine.Name,
			"quantity": item.Quantity,
			"price":    item.Price,
			"subtotal": item.Price * float64(item.Quantity),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"receipt": gin.H{
			"receipt_id": sale.ID,
			"cashier":    sale.User.Name,
			"customer":   sale.Customer.Name,
			"total":      sale.Total,
			"date":       sale.CreatedAt,
			"items":      items,
		},
	})
}

func DownloadReceiptPDF(c *gin.Context) {

	id := c.Param("id")

	var sale models.Sale

	err := config.DB.
		Preload("User").
		Preload("Customer").
		Preload("SaleItems").
		Preload("SaleItems.Medicine").
		First(&sale, id).Error

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Sale not found",
		})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, "Pharmacy Receipt")

	pdf.Ln(15)

	// Receipt Info
	pdf.SetFont("Arial", "", 12)

	pdf.Cell(40, 10, "Receipt ID:")
	pdf.Cell(40, 10, strconv.Itoa(int(sale.ID)))

	pdf.Ln(8)

	pdf.Cell(40, 10, "Customer:")
	pdf.Cell(40, 10, sale.Customer.Name)

	pdf.Ln(8)

	pdf.Cell(40, 10, "Cashier:")
	pdf.Cell(40, 10, sale.User.Name)

	pdf.Ln(15)

	// Table Header
	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(70, 10, "Medicine")
	pdf.Cell(30, 10, "Qty")
	pdf.Cell(40, 10, "Price")
	pdf.Cell(40, 10, "Subtotal")

	pdf.Ln(10)

	// Table Content
	pdf.SetFont("Arial", "", 12)

	for _, item := range sale.SaleItems {

		subtotal := item.Price * float64(item.Quantity)

		pdf.Cell(70, 10, item.Medicine.Name)
		pdf.Cell(30, 10, strconv.Itoa(item.Quantity))
		pdf.Cell(40, 10, strconv.FormatFloat(item.Price, 'f', 2, 64))
		pdf.Cell(40, 10, strconv.FormatFloat(subtotal, 'f', 2, 64))

		pdf.Ln(10)
	}

	pdf.Ln(10)

	// Total
	pdf.SetFont("Arial", "B", 14)

	pdf.Cell(40, 10, "Total:")

	pdf.Cell(
		40,
		10,
		strconv.FormatFloat(sale.Total, 'f', 2, 64),
	)

	fileName := "receipt.pdf"

	err = pdf.OutputFileAndClose(fileName)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate PDF",
		})
		return
	}

	c.File(fileName)
}
