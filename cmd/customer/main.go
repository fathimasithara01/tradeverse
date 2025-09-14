package main

import (
	"log"

	repositoryy "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/internal/customer/router"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"

	servicer "github.com/fathimasithara01/tradeverse/internal/admin/service"

	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go"
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

	userRepo := repositoryy.NewUserRepository(gormDB)
	roleRepo := repositoryy.NewRoleRepository(gormDB)

	userService := servicer.NewUserService(userRepo, roleRepo, cfg.JWTSecret)

	authController := controllers.NewAuthController(userService)
	profileController := controllers.NewProfileController(userService)

	kycRepo := repository.NewKYCRepository(gormDB)
	kycSvc := service.NewKYCService(kycRepo)
	kycController := controllers.NewKYCController(kycSvc)

	walletRepo := repository.NewWalletRepository(gormDB)
	pgClient := paymentgateway.NewSimulatedPaymentClient()
	walletSvc := service.NewWalletService(walletRepo, pgClient)
	walletController := controllers.NewWalletController(walletSvc)

	subscriptionRepo := repository.NewSubscriptionRepository(gormDB)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	adminTraderSubscriptionPlanRepo := repository.NewAdminTraderSubscriptionPlanRepository(gormDB)
	adminTraderSubscriptionPlanService := service.NewAdminTraderSubscriptionPlanService(adminTraderSubscriptionPlanRepo)
	adminTraderSubscriptionPlanController := controllers.NewAdminTraderSubscriptionPlanController(adminTraderSubscriptionPlanService)

	// traderSubscriptionRepo := repository.NewTraderSubscriptionRepository(gormDB)
	// traderSubscriptionService := service.NewTraderSubscriptionService(traderSubscriptionRepo, userRepo, pgClient)
	// traderSubscriptionController := controllers.NewTraderSubscriptionController(traderSubscriptionService)

	r := router.SetupRouter(cfg, authController, profileController, kycController, walletController, subscriptionController, adminTraderSubscriptionPlanController)

	port := cfg.CustomerPort
	log.Printf("Customer API server starting on port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start customer server: %v", err)
	}
}
