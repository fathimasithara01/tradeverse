// File: adminBackend/main.go

package main

import (
	"log"
	"net/http"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/db"
	"github.com/fathimasithara01/tradeverse/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration from .env file
	config.LoadConfig()

	// 2. Connect to the database
	db.ConnectDatabase()
	// Consider running migrations here, e.g.,
	// db.DB.AutoMigrate(&models.Admin{}, &models.Trader{})

	// 3. Set Gin mode from config
	// gin.SetMode(config.AppConfig.GinMode)

	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*.html")

	// 6. Setup routes
	routes.AdminRoutes(r)
	// routes.UserRoutes(r)

	// 7. *** FIX: The root route should redirect to the login page, NOT render a file. ***
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/login")
	})

	// 8. Start the server
	port := config.AppConfig.Port
	log.Printf("Server starting on port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
