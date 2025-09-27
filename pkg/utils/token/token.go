package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// VerifyToken verifies the JWT token
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Replace with your actual JWT secret key
	var jwtSecret = []byte("YOUR_SUPER_SECRET_JWT_KEY") // IMPORTANT: Use an environment variable for this!

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenID(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)
	token, err := VerifyToken(tokenString)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id claim not found or not a float64")
	}

	return uint(idFloat), nil
}

// ExtractTokenRole extracts the user role from a verified JWT token
func ExtractTokenRole(c *gin.Context) (models.UserRole, error) {
	tokenString := ExtractToken(c)
	token, err := VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", errors.New("invalid token claims")
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role claim not found or not a string")
	}

	return models.UserRole(roleStr), nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ExtractToken(c)
		if tokenString == "" {
			// Fix: Use c.JSON and c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		token, err := VerifyToken(tokenString)
		if err != nil || !token.Valid {
			// Fix: Use c.JSON and c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// TraderAuthMiddleware checks if the authenticated user is a trader
func TraderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := ExtractTokenRole(c)
		if err != nil {
			// Fix: Use c.JSON and c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Cannot determine user role"})
			c.Abort()
			return
		}
		if role != models.RoleTrader {
			// Fix: Use c.JSON and c.Abort()
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Only traders can access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}
