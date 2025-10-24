package main

import (
	"log"
	"strings"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/trader/bootstrap"
)

func main() {
	// Load config first
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(" Failed to load config: %v", err)
	}

	// Initialize the application
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Get trader port from config
	port := config.AppConfig.Server.TraderPort

	// Ensure port starts with ":" (required by Gin)
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf(" Trader service running on %s", port)

	// Start the server
	if err := app.Engine().Run(port); err != nil {
		log.Fatalf(" Server stopped with error: %v", err)
	}
}
