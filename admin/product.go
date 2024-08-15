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
	listorder := c.Query("list_order")
	var product []responsemodels.Product
	//tx := database.DB.Find(&product)
	// tx := database.DB.Raw(`SELECT * FROM categories join products on categories.id=products.category_id and products.deleted_at IS NULL AND categories.deleted_at IS NULL`).Scan(&product)
	sql := `SELECT * FROM categories join products on categories.id=products.category_id and products.deleted_at IS NULL AND categories.deleted_at IS NULL`
	if listorder == "" || listorder == "ASC" {
		sql += ` ORDER BY products.id ASC`
	} else if listorder == "DSC" {
		sql += ` ORDER BY products.id DESC`
	}
	tx := database.DB.Raw(sql).Scan(&product)
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
			"error_code": http.StatusBadRequest,
		})
		return
	}
	if Product.Price != float64(int(Product.Price)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "price should not contain decimal places",
		})
		return
	}

	p := Product.Size
	fmt.Println("---", Product.Size)
	if p != "Medium" && p != "Small" && p != "Large" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown size",
		})
		return
	}

	fmt.Println("product.categoryname ", Product.CategoryName)
	var count int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE categories.category_name=? and categories.deleted_at is NULL`, Product.CategoryName).Scan(&count).Error
	if err != nil {
		fmt.Println("failed to execute query", err)
	}
	fmt.Println("count", count)
	if count != 0 {
		var categoryid uint
		database.DB.Raw(`SELECT id from categories where category_name = ?`, Product.CategoryName).Scan(&categoryid)
		//var product models.Product
		product := models.Product{
			CategoryID:  categoryid,
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

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Product Added Successfully"})

}
func ProductEdit(c *gin.Context) {
	fmt.Println("hello")
	productID := c.Param("id")
	fmt.Println(productID)
	fmt.Println("hi")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM products WHERE id = ? AND deleted_at IS NULL`, productID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product id does not exist",
		})
		return
	}
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
	p := Product.Size
	fmt.Println("---", Product.Size)
	if p != "Medium" && p != "Small" && p != "Large" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown size",
		})
		return
	}
	var count1 int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM categories join products on categories.id=products.category_id WHERE categories.category_name=? and categories.deleted_at is NULL`, Product.CategoryName).Scan(&count1).Error
	if err != nil {
		fmt.Println("failed to execute query", err)
	}
	if count1 != 0 {
		var categoryid uint
		database.DB.Raw(`SELECT id from categories where category_name = ?`, Product.CategoryName).Scan(&categoryid)
		//var product models.Product
		product := models.Product{
			CategoryID:  categoryid,
			ProductName: Product.ProductName,
			Description: Product.Description,
			ImageUrl:    Product.ImageUrl,
			Price:       Product.Price,
			Stock:       Product.Stock,
			Popular:     Product.Popular,
			Size:        Product.Size,
		}
		database.DB.Model(&models.Product{}).Where("id = ?", productID).Updates(&product)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product updated successfully"})
}

func ProductDelete(c *gin.Context) {
	ProductID := c.Param("id")
	fmt.Println(ProductID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM products WHERE id = ? AND deleted_at IS NULL`, ProductID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product id does not exist",
		})
		return
	}

	database.DB.Where("id = ?", ProductID).Delete(&models.Product{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product deleted successfully"})
}
