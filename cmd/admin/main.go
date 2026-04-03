package main

import (
	"context"
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/bootstrap"
	"github.com/gin-contrib/cors"
)

func main() {

	ctx :=context.Background()
	// Initialize the application (DB, services, router)
	app, err := bootstrap.InitializeApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Apply CORS middleware
	app.Engine().Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change to specific domains in production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Run the server
	if err := app.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}