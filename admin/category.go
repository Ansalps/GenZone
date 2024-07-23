package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func Category(c *gin.Context) {
	var category []models.Category
	tx := database.DB.Find(&category)
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
			"categories": category,
		},
	})
}

func CategoryAdd(c *gin.Context) {
	fmt.Println("hello")
	var Category models.CategoryEdit
	err := c.BindJSON(&Category)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(Category); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	category := models.Category{
		CategoryName: Category.CategoryName,
		Description:  Category.Description,
		ImageUrl:     Category.ImageUrl,
	}
	database.DB.Create(&category)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})

}

func CategoryEdit(c *gin.Context) {
	fmt.Println("hello")
	CategoryID := c.Param("id")
	fmt.Println(CategoryID)
	var Category models.CategoryEdit
	err := c.BindJSON(&Category)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := helper.Validate(Category); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	category := models.Category{
		CategoryName: Category.CategoryName,
		Description:  Category.Description,
		ImageUrl:     Category.ImageUrl,
	}
	database.DB.Model(&models.Category{}).Where("id = ?", CategoryID).Updates(&category)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
}

func CategoryDelete(c *gin.Context) {
	fmt.Println("hello")
	CategoryID := c.Param("id")
	fmt.Println(CategoryID)
	// var Category models.Category
	// err := c.BindJSON(&Category)
	// response := gin.H{
	// 	"status":  false,
	// 	"message": "failed to bind request",
	// }
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, response)
	// 	return
	// }
	database.DB.Where("id = ?", CategoryID).Delete(&models.Category{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "category deleted succesfully"})
}
