package main

import (
	"log"

	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	adminService "github.com/fathimasithara01/tradeverse/internal/admin/service"
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

	userRepo := adminRepo.NewUserRepository(gormDB)
	roleRepo := adminRepo.NewRoleRepository(gormDB)
	userService := adminService.NewUserService(userRepo, roleRepo, cfg.JWTSecret)

	authController := controllers.NewAuthController(userService)

	tradeRepo := repository.NewTradeRepository(gormDB)
	subRepo := repository.NewSubscriberRepository(gormDB)
	liveRepo := repository.NewLiveTradeRepository(gormDB)

	tradeService := service.NewTradeService(tradeRepo)
	subService := service.NewSubscriberService(subRepo)
	liveService := service.NewLiveTradeService(liveRepo)

	tradeController := controllers.NewTradeController(tradeService)
	subController := controllers.NewSubscriberController(subService)
	liveController := controllers.NewLiveTradeController(liveService)

	r := router.SetupRouter(cfg, authController, tradeController, subController, liveController)

	port := cfg.TraderPort
	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start customer server: %v", err)
	}
}
