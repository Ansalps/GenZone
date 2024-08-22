package admin

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	listorder := c.Query("list_order")
	var users []responsemodels.User
	//tx := database.DB.Find(&users)
	sql := `SELECT * FROM users`
	if listorder == "" || listorder == "ASC" {
		sql += ` ORDER BY users.id ASC`
	} else if listorder == "DSC" {
		sql += ` ORDER BY users.id DESC`
	}
	tx := database.DB.Raw(sql).Scan(&users)
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
			"users": users,
		},
	})
}
