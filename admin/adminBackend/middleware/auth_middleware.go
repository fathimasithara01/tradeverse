package middleware

import (
	"log"
	"net/http"

	"github.com/fathimasithara01/tradeverse/auth"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("admin_token")
		if err != nil {
			log.Println("[MIDDLEWARE-ERROR] Cookie 'admin_token' not found.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			return
		}

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("[MIDDLEWARE-ERROR] Token validation failed: %s\n", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.Printf("[MIDDLEWARE-SUCCESS] Token validated. Claims: UserID=%d, Role=%s\n", claims.UserID, claims.Role)

		if claims.Role != "admin" {
			log.Printf("[MIDDLEWARE-FORBIDDEN] Access denied. UserID %d has role '%s', not 'admin'.\n", claims.UserID, claims.Role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: admin privileges required"})
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
