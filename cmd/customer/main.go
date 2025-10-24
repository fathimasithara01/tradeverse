package main

import (
	"log"
	"strings"
	"time"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/customer/bootstrap"
	"github.com/gin-contrib/cors"
)

func main() {
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	engine := app.Engine()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // You can replace "*" with specific domains for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	engine.Use(cors.New(corsConfig))

	port := config.AppConfig.Server.CustomerPort
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	if err := app.Engine().Run(port); err != nil {
		log.Fatalf(" Server stopped with error: %v", err)
	}
}
