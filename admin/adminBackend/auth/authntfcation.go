package auth

import (
	"time"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, email string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &AuthClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func ValidateJWT(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	secret := config.AppConfig.JWTSecret

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
