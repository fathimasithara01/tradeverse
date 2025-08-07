package main

import (
	"html/template"
	"log"

	"github.com/fathimasithara01/tradeverse/admin/config"
	database "github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	database.ConnectDatabase()
	// migration.AutoMigrations()
	log.Println(" Connected to PostgreSQL")

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{})
	// r.LoadHTMLGlob("adminDashboard/static/index.html")
	r.Static("/static", "../../adminDashboard/static")

	routes.AdminRoutes(r)
	routes.UserRoutes(r)

	port := config.AppConfig.Port
	if err := r.Run(":" + port); err != nil {
		log.Fatalf(" Failed to start server on port %s: %v", port, err)
	}

	log.Println(" Server running on port", port)
}
