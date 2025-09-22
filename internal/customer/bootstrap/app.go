package bootstrap

import (
	"log"

	adminSvc "github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"

	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"

	"github.com/fathimasithara01/tradeverse/internal/customer/router"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	port   string
}

func InitializeApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := migrations.ConnectDB(*cfg)
	if err != nil {
		return nil, err
	}

	userRepo := adminRepo.NewUserRepository(db)
	roleRepo := adminRepo.NewRoleRepository(db)
	kycRepo := customerrepo.NewKYCRepository(db)
	walletRepo := walletrepo.NewWalletRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	traderRepo := customerrepo.NewTraderRepository(db)
	subRepo := customerrepo.NewSubscriptionRepository(db)
	traderWalletRepo := walletrepo.NewTraderWalletRepository(db)

	userService := adminSvc.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	kycService := service.NewKYCService(kycRepo)
	paymentClient := paymentgateway.NewSimulatedPaymentClient()
	walletService := service.NewWalletService(walletRepo, paymentClient, db)
	customerService := service.NewCustomerService(customerRepo, walletService, walletRepo, db)
	traderService := service.NewTraderService(traderRepo, db)
	subService := service.NewSubscriptionService(db, subRepo)
	traderWalletService := service.NewTraderWalletService(db, traderWalletRepo)

	authController := controllers.NewAuthController(userService)
	profileController := controllers.NewProfileController(userService)
	kycController := controllers.NewKYCController(kycService)
	walletController := controllers.NewWalletController(walletService)
	customerController := controllers.NewCustomerController(customerService)
	traderController := controllers.NewTraderController(traderService)
	subController := controllers.NewSubscriptionController(subService)
	traderWalletController := controllers.NewTraderWalletController(traderWalletService)

	r := router.SetupRouter(
		cfg,
		authController,
		profileController,
		kycController,
		walletController,
		customerController,
		traderController,
		subController,
		traderWalletController,
	)

	return &App{
		engine: r,
		port:   cfg.CustomerPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Customer API server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
