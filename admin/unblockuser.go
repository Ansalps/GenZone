package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func UnblockUser(c *gin.Context) {
	unblockID := c.Param("id")
	fmt.Println(unblockID)
	database.DB.Model(&models.User{}).Where("id = ?", unblockID).Update("Status", "Active")
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "user unblocked"})
}
