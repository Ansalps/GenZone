package public

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func ReadCategory(c *gin.Context) {
	var category []responsemodels.Category
	sql := `SELECT * FROM categories`

	tx := database.DB.Raw(sql).Scan(&category)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to retrieve data from the database, or the data doesn't exists",
		})
		return
	}
	for i := range category {
		if category[i].ImageUrl != "" {
			if resolvedURL, err := helper.ResolveObjectURL(c.Request.Context(), "S3_BUCKET_NAME", category[i].ImageUrl); err == nil {
				category[i].ImageUrl = resolvedURL
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"categories": category,
		},
	})
}
