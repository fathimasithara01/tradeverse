package middleware

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/auth"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("admin_token")

		// If the cookie is not set, redirect to login page.
		if err != nil {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// Validate the token.
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			// If token is invalid, clear the bad cookie and redirect.
			c.SetCookie("admin_token", "", -1, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// Set admin details in the context for use in controllers.
		c.Set("admin_id", claims.ID)
		c.Set("admin_email", claims.Email)
		c.Set("admin_role", claims.Role)

		c.Next()
	}
}
