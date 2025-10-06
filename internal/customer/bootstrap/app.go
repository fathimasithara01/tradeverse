package bootstrap

import (
	"log"

	adminSvc "github.com/fathimasithara01/tradeverse/internal/admin/service"

	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	traderSignalRepo "github.com/fathimasithara01/tradeverse/internal/trader/repository"

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
	walletRepo := walletrepo.NewWalletRepository(db)
	kycRepo := customerrepo.NewKYCRepository(db)
	// subRepo := customerrepo.NewSubscriptionRepository(db)
	traderRepo := customerrepo.NewTraderRepository(db)
	adminSubRepo := customerrepo.NewIAdminSubscriptionRepository(db)
	// traderWalletRepo := walletrepo.NewTraderWalletRepository(db)
	customerTraderSubRepo := customerrepo.NewTraderSubscriptionRepository(db)
	traderSignalRepo := traderSignalRepo.NewSignalRepository(db)

	userService := adminSvc.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	kycService := service.NewKYCService(kycRepo)
	paymentClient := paymentgateway.NewSimulatedPaymentClient()
	walletService := service.NewWalletService(walletRepo, paymentClient, db)
	traderService := service.NewTraderService(traderRepo, db)
	// subService := service.NewSubscriptionService(db, subRepo)
	// traderWalletService := service.NewTraderWalletService(db, traderWalletRepo)
	adminSubService := service.NewCustomerService(adminSubRepo, walletService, walletRepo, db)
	customerTraderSubSvc := service.NewTraderSubscriptionService(customerTraderSubRepo, db)
	customerSignalSvc := service.NewCustomerSignalService(traderSignalRepo, customerTraderSubSvc)

	authController := controllers.NewAuthController(userService)
	profileController := controllers.NewProfileController(userService)
	kycController := controllers.NewKYCController(kycService)
	walletController := controllers.NewWalletController(walletService)
	adminSubController := controllers.NewAdminSubscriptionController(adminSubService)
	traderController := controllers.NewTraderController(traderService)
	customerTraderSubCtrl := controllers.NewTraderSubscriptionController(customerTraderSubSvc)
	customerSignalCtrl := controllers.NewCustomerSignalController(customerSignalSvc)

	// subController := controllers.NewSubscriptionController(subService)
	// traderWalletController := controllers.NewTraderWalletController(traderWalletService)

	r := router.SetupRouter(
		cfg,
		authController,
		profileController,
		kycController,
		walletController,
		adminSubController,
		traderController,
		customerTraderSubCtrl,
		customerSignalCtrl,
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
