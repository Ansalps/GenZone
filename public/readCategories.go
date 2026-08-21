package public

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func ReadCategory(c *gin.Context) {

	
	var category []responsemodels.Category
	//tx := database.DB.Find(&category)
	sql := `SELECT * FROM categories`

	tx := database.DB.Raw(sql).Scan(&category)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
