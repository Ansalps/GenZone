package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Product(c *gin.Context) {
	var product []responsemodels.Product
	//tx := database.DB.Find(&product)
	tx := database.DB.Raw(`SELECT * FROM categories join products on categories.id=products.category_id and products.deleted_at IS NULL`).Scan(&product)
	if tx.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "failed to retrieve data from the database, or the data doesn't exists",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"products": product,
		},
	})

}
func ProductAdd(c *gin.Context) {
	fmt.Println("hello")
	var Product models.ProductAdd
	err := c.BindJSON(&Product)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := helper.Validate(Product); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE id=? and deleted_at is NULL`, Product.CategoryID).Scan(&count).Error
	if err != nil {
		fmt.Println("failed to execute query", err)
	}
	if count != 0 {
		//var product models.Product
		product := models.Product{
			CategoryID:  Product.CategoryID,
			ProductName: Product.ProductName,
			Description: Product.Description,
			ImageUrl:    Product.ImageUrl,
			Price:       Product.Price,
			Stock:       Product.Stock,
			Popular:     Product.Popular,
			Size:        Product.Size,
		}
		database.DB.Create(&product)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})

}
func ProductEdit(c *gin.Context) {
	fmt.Println("hello")
	productID := c.Param("id")
	fmt.Println(productID)
	fmt.Println("hi")
	var Product models.ProductEdit
	err := c.BindJSON(&Product)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := helper.Validate(Product); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	database.DB.Model(&models.Product{}).Where("id = ?", productID).Updates(&Product)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
}

func ProductDelete(c *gin.Context) {
	fmt.Println("hello")
	ProductID := c.Param("id")
	fmt.Println("hello")
	fmt.Println(ProductID)
	fmt.Println("hi")

	database.DB.Where("id = ?", ProductID).Delete(&models.Product{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product deleted successfully"})
}
