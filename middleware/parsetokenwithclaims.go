package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CustomClaims struct

// Secret key
var Secret = []byte("your-secret-key")

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the token from the cookie that matches the required role.
		cookieName := "jwt_user_token"
		if requiredRole == "admin" {
			cookieName = "jwt_admin_token"
		}

		tokenString, err := c.Cookie(cookieName)
		if err != nil {
			// Backward compatibility for older single-cookie sessions.
			tokenString, err = c.Cookie("jwt_token")
			if err != nil {
				fmt.Println("error", err)
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Please Log In"})
				c.Abort()
				return
			}
		}

		claims := &CustomClaims{}
		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			if claims.Role != requiredRole {
				c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient privileges"})
				c.Abort()
				return
			}
			// Store claims in context for further use in handlers
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
