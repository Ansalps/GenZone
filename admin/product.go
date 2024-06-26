package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func Product(c *gin.Context) {
	var cy models.Product
	database.DB.Find(&cy)
	fmt.Println(cy)

}
func ProductAdd(c *gin.Context) {
	fmt.Println("hello")
	var Product models.Product
	err := c.BindJSON(&Product)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var count int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE id=?`, Product.CategoryID).Scan(&count).Error
	if err != nil {
		database.DB.Create(&Product)
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
	var Product models.Product
	err := c.BindJSON(&Product)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
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

}
