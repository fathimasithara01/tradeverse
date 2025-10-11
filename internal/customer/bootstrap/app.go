package bootstrap

import (
	"log"

	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository" // Use adminRepo alias consistently
	adminSvc "github.com/fathimasithara01/tradeverse/internal/admin/service"
	traderControllers "github.com/fathimasithara01/tradeverse/internal/trader/controllers"
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	traderSvc "github.com/fathimasithara01/tradeverse/internal/trader/service"

	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"

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

	userRepo := adminRepo.NewUserRepository(db) // admin user repo
	roleRepo := adminRepo.NewRoleRepository(db)
	// adminWalletRepo := adminRepo.NewAdminWalletRepository(db) // Renamed for clarity
	kycRepo := customerrepo.NewKYCRepository(db)
	traderRepo := customerrepo.NewTraderRepository(db)
	adminSubRepo := customerrepo.NewIAdminSubscriptionRepository(db) // Typo: should be NewAdminSubscriptionRepository?
	customerWalletRepo := walletrepo.NewWalletRepository(db)
	// customerSignalRepo := customerrepo.NewCustomerSignalRepository(db)
	// customerUserRepo := customerrepo.NewUserRepository(db) // Renamed for clarity vs admin userRepo
	// subscriptionPlanRepo := customerrepo.NewSubscriptionPlanRepository(db)
	customerTraderSubsRepo := customerrepo.NewCustomerTraderSubscriptionRepository(db)
	traderSubsRepo := repository.NewTraderSubscriptionRepository(db)

	// --- Services ---
	userService := adminSvc.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	kycService := service.NewKYCService(kycRepo)
	paymentClient := paymentgateway.NewSimulatedPaymentClient()
	walletService := service.NewWalletService(db, customerWalletRepo, paymentClient)
	traderService := service.NewTraderService(traderRepo, db)
	adminSubService := service.NewAdminSubscriptionService(adminSubRepo, walletService, customerWalletRepo, db) // Assuming this is the correct constructor
	customerTraderSubsService := service.NewCustomerTraderSubscriptionService(customerTraderSubsRepo, db)
	traderSubsService := traderSvc.NewTraderSubscriptionService(traderSubsRepo, db) // Pass db for transactions

	customerTraderSubsController := controllers.NewCustomerTraderSubscriptionController(customerTraderSubsService)
	// customerSubscriptionController := controllers.NewCustomerSubscriptionController(traderSubscriptionService, customerSignalService)
	// customerSignalController := controllers.NewCustomerSignalController(customerSignalService) // This now correctly injects signalService
	authController := controllers.NewAuthController(userService)
	profileController := controllers.NewProfileController(userService)
	kycController := controllers.NewKYCController(kycService)
	walletController := controllers.NewWalletController(walletService)
	adminSubController := controllers.NewAdminSubscriptionController(adminSubService)
	traderController := controllers.NewTraderController(traderService)
	traderSubsController := traderControllers.NewTraderSubscriptionController(traderSubsService)

	// Unused controllers (commented out in your original)
	// customerTraderSubCtrl := controllers.NewTraderSubscriptionController(customerTraderSubSvc) // Requires customerTraderSubSvc if exists
	// customerSignalCtrl := controllers.NewCustomerSignalController(customerSignalSvc) // Requires customerSignalSvc if exists
	// subController := controllers.NewSubscriptionController(subService) // Requires subService if exists
	// traderWalletController := controllers.NewTraderWalletController(traderWalletService) // Requires traderWalletService if exists

	r := router.SetupRouter(
		cfg,
		authController,
		profileController,
		kycController,
		walletController,
		adminSubController,
		traderController,
		customerTraderSubsController,
		traderSubsController,
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
