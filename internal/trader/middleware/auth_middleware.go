package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TraderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role type"})
			return
		}

		if userRole != "trader" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: Traders only"})
			return
		}

		c.Next()
	}
}
