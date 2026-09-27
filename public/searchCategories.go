package public

import (
	"net/http"
	"strings"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	search := c.Query("search")
	nameSort := c.Query("name_sort")

	newArrivals := c.Query("new_arrivals")

	var category []responsemodels.Category
	//tx := database.DB.Find(&category)
	sql := `SELECT * FROM categories`

	var args []interface{}

	// --------------------------------
	// Search Filter
	// --------------------------------
	if search != "" {
		sql += ` WHERE category_name ILIKE ?`
		args = append(args, "%"+search+"%")
	}

	// --------------------------------
	// Sorting
	// --------------------------------
	var orderBy []string

	if nameSort != "" {
		switch nameSort {
		case "aA-zZ":
			orderBy = append(orderBy, "category_name ASC")
		case "zZ-aA":
			orderBy = append(orderBy, "category_name DESC")
		default:
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid name sort filter"})
			return
		}
	}

	if newArrivals != "" {
		switch newArrivals {
		case "true":
			orderBy = append(orderBy, "created_at DESC")
		case "false":
			// No sorting required
		default:
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "new_arrivals must be true or false"})
			return
		}
	}

	if len(orderBy) == 0 {
		orderBy = append(orderBy, "created_at DESC")
	}

	sql += ` ORDER BY ` + strings.Join(orderBy, ", ")

	tx := database.DB.Raw(sql,args...).Scan(&category)
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
