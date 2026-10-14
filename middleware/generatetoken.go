package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

// Function to create a JWT token
func CreateToken(role string, email string, id uint) (string, error) {
	claims := CustomClaims{
		Email: email,
		Role:  role,
		ID:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "GenZone",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
