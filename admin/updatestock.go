package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockUpdateItem struct {
	Size  string `json:"size" binding:"required"`
	Stock uint   `json:"stock"`
}

type UpdateStockRequest struct {
	Inventory []StockUpdateItem `json:"inventory" binding:"required"`
}

func UpdateProductStock(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Product id is required"})
		return
	}

	var req UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid stock payload"})
		return
	}

	if len(req.Inventory) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Inventory cannot be empty"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to load product"})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to start transaction"})
		return
	}

	for _, item := range req.Inventory {
		size := strings.TrimSpace(item.Size)
		if size == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Size is required for each stock entry"})
			return
		}

		var variant models.ProductVariant
		err := tx.Where("product_id = ? AND size = ?", product.ID, size).First(&variant).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newVariant := models.ProductVariant{
				ProductID: product.ID,
				Size:      size,
				Stock:     item.Stock,
			}
			if err := tx.Create(&newVariant).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create product variant stock"})
				return
			}
			continue
		}
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to fetch product variant"})
			return
		}

		if err := tx.Model(&variant).Update("stock", item.Stock).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update product stock"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to commit stock update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Stock updated successfully"})
}
