package server

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/intenal/graph"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil || token != "dummy-secret-token" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("userID", "user-123")
		c.Next()
	}
}

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/login", func(c *gin.Context) { c.HTML(200, "portal.html", gin.H{"page": "login"}) })
	r.POST("/login", func(c *gin.Context) {
		c.SetCookie("auth_token", "dummy-secret-token", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
	})
	app := r.Group("/")
	app.Use(AuthMiddleware())
	{
		app.GET("/", func(c *gin.Context) { c.HTML(200, "portal.html", gin.H{"page": "app"}) })
		gqlHandler := handler.New(&handler.Config{Schema: &graph.Schema, Pretty: true})
		app.POST("/graphql", gin.WrapH(gqlHandler))
		app.GET("/graphql", gin.WrapH(gqlHandler))
	}
	return r
}
