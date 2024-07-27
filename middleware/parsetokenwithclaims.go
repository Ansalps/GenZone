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
	fmt.Println("hi")
	return func(c *gin.Context) {
		// Get token from cookie
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please Log In"})
			c.Abort()
			return
		}
		//Authorization cookie required
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check user role
		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"message": "Log in to continue"})
			c.Abort()
			return
		}
		//Insufficient privileges

		// Set claims in context
		//c.Set("claims", claims)
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Next()
	}
}
