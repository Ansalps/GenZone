package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CustomClaims struct

// Secret key
var Secret = []byte("your-secret-key")

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	fmt.Println("hi")
	return func(c *gin.Context) {
		//Get token from cookie
		// tokenString, err := c.Cookie("jwt_token")
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Please Log In"})
		// 	c.Abort()
		// 	return
		// }
		// //Authorization cookie required...
		// claims := &CustomClaims{}
		// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return Secret, nil
		// })

		// // Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Println("---", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("-----------------------", tokenString)
		//claims := &CustomClaims{}
		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		var claim CustomClaims
		fmt.Println("ggg", claim.Role)
		// Check user role
		//fmt.Println("fff",&CustomClaims.Role)

		//Insufficient privileges

		// Set claims in context
		//c.Set("claims", claims)
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
		// fmt.Println("role1", claims.Role)
		// fmt.Println("role2", requiredRole)
		// if claims.Role != requiredRole {
		// 	c.JSON(http.StatusForbidden, gin.H{"message": "Log in to continue"})
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
