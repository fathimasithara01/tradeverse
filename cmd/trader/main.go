package main

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/trader/bootstrap"
)

func main() {
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
