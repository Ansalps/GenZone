package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var AdminLogin models.AdminLogin
	err := c.BindJSON(&AdminLogin)

	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
	fmt.Println("hi")
}

func Category(c *gin.Context) {
	fmt.Println("hi")
	var AdminCategory models.AdminCategory
	err := c.BindJSON(&AdminCategory)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	database.DB.Find(&models.AdminCategory{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
}

func CategoryAdd(c *gin.Context) {
	fmt.Println("hello")
	var AdminCategory models.AdminCategory
	err := c.BindJSON(&AdminCategory)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	database.DB.Create(&AdminCategory)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})

}

func CategoryEdit(c *gin.Context) {
	fmt.Println("hello")
	CategoryID := c.Param("id")
	fmt.Println(CategoryID)
	var AdminCategory models.AdminCategory
	err := c.BindJSON(&AdminCategory)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	database.DB.Model(&models.AdminCategory{}).Where("id = ?", CategoryID).Updates(&AdminCategory)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
}

func CategoryDelete(c *gin.Context) {
	fmt.Println("hello")
	CategoryID := c.Param("id")
	fmt.Println(CategoryID)
	var AdminCategory models.AdminCategory
	err := c.BindJSON(&AdminCategory)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	database.DB.Where("id = ?", CategoryID).Delete(&models.AdminCategory{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
}
