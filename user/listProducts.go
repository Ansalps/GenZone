package user

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func ListProducts(c *gin.Context) {
	var Products []responsemodels.Product
	//tx := database.DB.Select("*").Find(&Products)
	tx := database.DB.Raw(`SELECT * FROM categories JOIN products ON products.category_id=categories.id WHERE products.deleted_at IS NULL AND categories.deleted_at IS NULL ORDER BY products.id`).Scan(&Products)
	if tx.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     false,
			"message":    "failed to retrieve data from the database, or the product doesn't exist",
			"error_code": http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved products",
		"data": gin.H{
			"products": Products,
		},
	})

}
