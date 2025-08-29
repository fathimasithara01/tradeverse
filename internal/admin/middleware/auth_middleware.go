package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/fathimasithara01/tradeverse/pkg/auth"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/service"
	"github.com/gin-gonic/gin"
)

type AuthzMiddleware struct {
	RoleSvc service.IRoleService
}

func NewAuthzMiddleware(roleSvc service.IRoleService) *AuthzMiddleware {
	return &AuthzMiddleware{RoleSvc: roleSvc}
}

func (a *AuthzMiddleware) RequirePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, roleExists := c.Get("userRole")
		roleID, idExists := c.Get("roleID")

		if !roleExists || !idExists {
			log.Println("[AUTHZ-ERROR] Role or RoleID not found in token context. Check JWTMiddleware.")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: Incomplete user token."})
			return
		}

		if userRole.(string) == "admin" {
			log.Printf("[AUTHZ-INFO] User is an admin. Bypassing permission check for '%s'. Access granted.", permissionName)
			c.Next()
			return
		}

		log.Printf("[AUTHZ-INFO] User is not an admin. Checking for permission '%s' for RoleID %v.", permissionName, roleID)
		hasPerm, err := a.RoleSvc.RoleHasPermission(roleID.(uint), permissionName)
		if err != nil || !hasPerm {
			log.Printf("[AUTHZ-DENIED] Access denied for RoleID %v, permission '%s'.", roleID, permissionName)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: You do not have the required permission."})
			return
		}

		log.Printf("[AUTHZ-SUCCESS] Access granted for RoleID %v, permission '%s'.", roleID, permissionName)
		c.Next()
	}
}

func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		cookie, err := c.Cookie("admin_token")
		if err == nil {
			tokenString = cookie
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				log.Println("[AUTH-WARN] No token cookie and no Authorization header found.")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
				return
			}
			tokenString = parts[1]
		}

		secret := cfg.JWTSecret
		claims, err := auth.ValidateJWT(tokenString, secret)
		if err != nil {
			log.Printf("[AUTH-ERROR] Token validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Set("roleID", claims.RoleID)

		c.Next()
	}
}
