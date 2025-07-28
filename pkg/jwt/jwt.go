package jwt

import (
	"fmt"
	"shopify-app/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID utils.BinaryUUID `json:"user_id"`
	Email  string           `json:"email"`
	Role   string           `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID utils.BinaryUUID, email, role, jwtSecret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(tokenString, jwtSecret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}