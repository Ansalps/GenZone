package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie(
		"jwt_user_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)
	c.SetCookie(
		"jwt_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Logged out successfully",
	})
}
