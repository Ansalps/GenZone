package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func BlockUser(c *gin.Context) {
	blockID := c.Param("id")
	fmt.Println(blockID)
	database.DB.Model(&models.User{}).Where("id = ?", blockID).Update("Status", "Blocked")
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "user blocked"})
}
