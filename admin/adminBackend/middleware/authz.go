// middleware/authz.go

package middleware

import (
	"fmt"
	"net/http"

	"github.com/fathimasithara01/tradeverse/models" // Import models for the RoleAdmin constant
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type AuthzMiddleware struct {
	RoleSvc *service.RoleService
}

func NewAuthzMiddleware(roleSvc *service.RoleService) *AuthzMiddleware {
	return &AuthzMiddleware{RoleSvc: roleSvc}
}

func (am *AuthzMiddleware) RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleVal, _ := c.Get("role")
		if role, ok := roleVal.(string); ok && role == string(models.RoleAdmin) {
			fmt.Printf("[AUTHZ-BYPASS] User is an admin. Granting automatic access to '%s'.\n", requiredPermission)
			c.Next()
			return // Exit immediately after granting access.
		}

		// --- REGULAR PERMISSION CHECK FOR NON-ADMINS ---
		roleIDVal, exists := c.Get("roleID")
		if !exists || roleIDVal.(uint) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: User has an invalid or unassigned Role ID."})
			return
		}
		roleID := roleIDVal.(uint)

		fmt.Printf("[AUTHZ-CHECK] RoleID: %d, Required Permission: '%s'\n", roleID, requiredPermission)

		hasPerm, err := am.RoleSvc.RoleHasPermission(roleID, requiredPermission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not verify permissions."})
			return
		}

		if !hasPerm {
			fmt.Printf("[ACCESS-DENIED] RoleID %d was denied access to '%s'.\n", roleID, requiredPermission)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: You do not have the required permission."})
			return
		}

		fmt.Printf("[ACCESS-GRANTED] RoleID %d was granted access to '%s'.\n", roleID, requiredPermission)
		c.Next()
	}
}
