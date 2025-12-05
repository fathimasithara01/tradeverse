package main

import (
	"log"
	"strings"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/trader/bootstrap"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(" Failed to load config: %v", err)
	}

	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	port := config.AppConfig.Server.TraderPort

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	if err := app.Engine().Run(port); err != nil {
		log.Fatalf(" Server stopped with error: %v", err)
	}
}
