package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func Product(c *gin.Context) {
	fmt.Println("hi")
	var AdminProduct models.AdminProduct
	err := c.BindJSON(&AdminProduct)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	database.DB.Find(&models.AdminProduct{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "done"})
}
