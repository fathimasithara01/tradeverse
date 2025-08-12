package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT Claims
type AuthClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new token.
func GenerateJWT(userID uint, email, role string) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is missing from config")
	}

	claims := AuthClaims{
		ID:    userID,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT checks a token string.
func ValidateJWT(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	secret := config.AppConfig.JWTSecret

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// --- Password Hashing ---

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// import (
// 	"fmt"
// 	"time"

// 	"github.com/fathimasithara01/tradeverse/config"
// 	"github.com/golang-jwt/jwt/v5"
// )

// func GenerateJWT(userID uint, email, role string) (string, error) {
// 	secret := config.AppConfig.JWTSecret
// 	if secret == "" {
// 		return "", fmt.Errorf("JWT_SECRET is missing")
// 	}

// 	claims := jwt.MapClaims{
// 		"id":    userID,
// 		"email": email,
// 		"role":  role,
// 		"exp":   time.Now().Add(24 * time.Hour).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(secret))
// }
// import "golang.org/x/crypto/bcrypt"

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	re
