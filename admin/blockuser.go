package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func BlockUser(c *gin.Context) {
	//blockID := c.Param("id")
	//fmt.Println(blockID)
	var blockuser models.BlockUser
	err := c.BindJSON(&blockuser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request failed to bind",
		})
		return
	}
	if err := helper.Validate(blockuser); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM users WHERE id = ? AND deleted_at IS NULL`, blockuser.UserID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user does not exist",
		})
		return
	}

	database.DB.Model(&models.User{}).Where("id = ?", blockuser.UserID).Update("Status", "Blocked")
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "user blocked"})
}
