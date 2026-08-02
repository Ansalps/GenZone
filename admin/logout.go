package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {

	//c.SetSameSite(http.SameSiteLaxMode) 
	// If your login has no explicit SetSameSite line, do NOT use SetSameSite here either.
	// Keeping it blank matches the "null" state perfectly so the browser allows the deletion.
	c.SetCookie(
		"jwt_token",
		"",
		-1, // -1 deletes it instantly
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
