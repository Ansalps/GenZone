package middleware

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Email string
	Role  string `json:"role"`
	ID    uint
	jwt.StandardClaims
}

// Function to create a JWT token
func CreateToken(role string, email string, id uint) (string, error) {
	claims := CustomClaims{
		Email: email,
		Role:  role,
		ID:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "GenZone",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret)
}
