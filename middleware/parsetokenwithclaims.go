package middleware

import (
	"net/http"
	"os"

	"github.com/Ansalps/GeZOne/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := "jwt_user_token"
		if requiredRole == "admin" {
			cookieName = "jwt_admin_token"
		}

		tokenString, err := c.Cookie(cookieName)
		if err != nil {

			tokenString, err = c.Cookie("jwt_token")
			if err != nil {

				c.JSON(http.StatusUnauthorized, gin.H{"message": "Please Log In"})
				c.Abort()
				return
			}
		}

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
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

			if requiredRole == "user" {
				var status string
				err := database.DB.Table("users").Select("status").Where("id = ?", claims.ID).Scan(&status).Error

				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "User account no longer exists"})
					c.Abort()
					return
				}

				if status == "Blocked" {
					c.JSON(http.StatusForbidden, gin.H{"message": "Your account has been suspended by an administrator."})
					c.Abort()
					return
				}
			}

			c.Set("claims", claims)
		} else {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
