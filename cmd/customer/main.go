package main

import (
	"log"

	repositoryy "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	servicer "github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository"

	customerrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	walletrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"

	"github.com/fathimasithara01/tradeverse/internal/customer/router"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"

	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
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

	kycRepo := customerrepo.NewKYCRepository(gormDB)
	kycSvc := service.NewKYCService(kycRepo)
	kycController := controllers.NewKYCController(kycSvc)

	paymentClient := paymentgateway.NewSimulatedPaymentClient()

	walletRepo := walletrepo.NewWalletRepository(gormDB)
	walletService := service.NewWalletService(walletRepo, paymentClient, gormDB)
	walletCtrl := controllers.NewWalletController(walletService)

	customerRepo := repository.NewCustomerRepository(gormDB)
	customerService := service.NewCustomerService(customerRepo, walletService, walletRepo, gormDB)
	customerController := controllers.NewCustomerController(customerService)

	traderRepo := customerrepo.NewTraderRepository(gormDB)
	traderService := service.NewTraderService(traderRepo, gormDB)
	traderController := controllers.NewTraderController(traderService)

	subRepo := customerrepo.NewSubscriptionRepository(gormDB)
	subSvc := service.NewSubscriptionService(gormDB, subRepo)
	subCntrl := controllers.NewSubscriptionController(subSvc)

	subTradRepo := walletrepo.NewTraderWalletRepository(gormDB)
	subTradeSvc := service.NewTraderWalletService(gormDB, subTradRepo)
	subTradeCntrl := controllers.NewTraderWalletController(subTradeSvc)

	r := router.SetupRouter(
		cfg,
		authController,
		profileController,
		kycController,
		walletCtrl,
		customerController,
		traderController,
		subCntrl,
		subTradeCntrl,
	)

	port := cfg.CustomerPort
	log.Printf("Customer API server starting on port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start customer server: %v", err)
	}
}
