package main

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/router"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/db"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
	"github.com/fathimasithara01/tradeverse/pkg/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	gormDB, err := db.ConnectDB(*cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := repository.NewUserRepository(gormDB)
	roleRepo := repository.NewRoleRepository(gormDB)

	userService := service.NewUserService(userRepo, roleRepo, cfg.JWTSecret)

	authController := controllers.NewAuthController(userService)
	profileController := controllers.NewProfileController(userService) // New

	r := router.SetupRouter(cfg, authController, profileController)

	customerPort := "8081"
	log.Printf("Customer API server starting on port http://localhost:%s", customerPort)
	if err := r.Run(":" + customerPort); err != nil {
		log.Fatalf("Failed to start customer server: %v", err)
	}
}
