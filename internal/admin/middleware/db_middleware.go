// pkg/middleware/db_middleware.go (or similar)
package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DBMiddleware injects the GORM DB instance into the Gin context.
func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db) // Set the DB instance in the context
		c.Next()        // Proceed to the next handler/middleware
	}
}