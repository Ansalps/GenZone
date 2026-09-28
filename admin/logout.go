package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {

	// Clear the dedicated admin token and the legacy fallback to avoid stale sessions.
	c.SetCookie(
		"jwt_admin_token",
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
