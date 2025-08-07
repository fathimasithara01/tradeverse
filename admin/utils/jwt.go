package utils

import (
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/admin/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, email, role string) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is missing")
	}

	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
