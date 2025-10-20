package bootstrap

import (
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	adminSvc "github.com/fathimasithara01/tradeverse/internal/admin/service"

	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"

	"github.com/fathimasithara01/tradeverse/internal/customer/router"
	"github.com/fathimasithara01/tradeverse/migrations"
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
	adminSubscriptionPlanRepo := adminRepo.NewSubscriptionPlanRepository(db)
	adminAdminWalletRepo := adminRepo.NewAdminWalletRepository(db)
	adminUserRepo := adminRepo.NewUserRepository(db)

	customerSubscriptionPlanRepo := customerrepo.NewCustomerSubscriptionPlanRepository(db)
	customerSubscriptionRepo := customerrepo.NewCustomerSubscriptionRepository(db)
	customerWalletRepo := walletrepo.NewWalletRepository(db)
	kycRepo := customerrepo.NewKYCRepository(db)
	traderRepo := customerrepo.NewTraderRepository(db)
	customerTraderSubsRepo := customerrepo.NewCustomerTraderSignalSubscriptionRepository(db)

	// --- Services ---
	adminAdminWalletService := adminSvc.NewAdminWalletService(adminAdminWalletRepo, db)
	customerWalletService := service.NewWalletService(db, customerWalletRepo, paymentgateway.NewSimulatedPaymentClient())
	customerSubscriptionPlanService := service.NewCustomerSubscriptionPlanService(customerSubscriptionPlanRepo)
	customerSubscriptionService := service.NewCustomerSubscriptionService(
		customerSubscriptionRepo,
		adminSubscriptionPlanRepo,
		adminAdminWalletService,
		adminUserRepo,
		db,
	)
	userService := adminSvc.NewUserService(userRepo, roleRepo, cfg.JWT.Secret)
	kycService := service.NewKYCService(kycRepo)
	paymentClient := paymentgateway.NewSimulatedPaymentClient()
	walletService := service.NewWalletService(db, customerWalletRepo, paymentClient)
	traderService := service.NewTraderService(traderRepo, db)
	customerTraderSubsService := service.NewCustomerTraderSignalSubscriptionService(customerTraderSubsRepo, db)

	subscriptionPlanController := controllers.NewSubscriptionPlanController(
		customerSubscriptionPlanService,
		customerSubscriptionService,
		customerWalletService,
	)

	// planService service.ICustomerSubscriptionPlanService,
	// subService service.ICustomerSubscriptionService,
	// walletService service.IWalletService,
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
		subscriptionPlanController,
	)

	return &App{
		engine: r,
		port:   cfg.Server.CustomerPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Customer API server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
