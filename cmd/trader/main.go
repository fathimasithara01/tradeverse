package main

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/trader/controllers"
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/internal/trader/router"
	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	gormDB, err := migrations.ConnectDB(*cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	tradeRepo := repository.NewTradeRepository(gormDB)
	walletRepo := repository.NewWalletRepository(gormDB)
	transactionRepo := repository.NewTransactionRepository(gormDB)

	tradeService := service.NewTradeService(tradeRepo, walletRepo, transactionRepo, gormDB)

	tradeController := controllers.NewTradeController(tradeService)

	r := router.SetupRouter(cfg, tradeController)

	port := cfg.TraderPort
	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start customer server: %v", err)
	}
}
