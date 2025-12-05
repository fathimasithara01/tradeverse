package main

import (
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/bootstrap"
	"github.com/gin-contrib/cors"
)

func main() {
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	app.Engine().Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if err := app.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
