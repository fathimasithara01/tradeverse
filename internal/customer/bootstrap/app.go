package bootstrap

import (
	"log"

	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	adminSvc "github.com/fathimasithara01/tradeverse/internal/admin/service"

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

	userRepo := adminRepo.NewUserRepository(db)
	roleRepo := adminRepo.NewRoleRepository(db)
	// adminWalletRepo := adminRepo.NewAdminWalletRepository(db)
	kycRepo := customerrepo.NewKYCRepository(db)
	traderRepo := customerrepo.NewTraderRepository(db)
	customerWalletRepo := walletrepo.NewWalletRepository(db)
	customerTraderSubsRepo := customerrepo.NewCustomerTraderSignalSubscriptionRepository(db)

	// --- Services ---
	userService := adminSvc.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	kycService := service.NewKYCService(kycRepo)
	paymentClient := paymentgateway.NewSimulatedPaymentClient()
	walletService := service.NewWalletService(db, customerWalletRepo, paymentClient)
	traderService := service.NewTraderService(traderRepo, db)
	customerTraderSubsService := service.NewCustomerTraderSignalSubscriptionService(customerTraderSubsRepo, db)

	customerTraderSubsController := controllers.NewCustomerTraderSignalSubscriptionController(customerTraderSubsService)
	authController := controllers.NewAuthController(userService)
	profileController := controllers.NewProfileController(userService)
	kycController := controllers.NewKYCController(kycService)
	walletController := controllers.NewWalletController(walletService)
	traderController := controllers.NewTraderController(traderService)

	r := router.SetupRouter(
		cfg,
		authController,
		profileController,
		kycController,
		walletController,
		traderController,
		customerTraderSubsController,
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
