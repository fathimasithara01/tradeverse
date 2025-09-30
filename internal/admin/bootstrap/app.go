package bootstrap

import (
	"fmt"
	"log"
	"path/filepath" // Import for path manipulation
	"runtime"       // Import for getting current file path
	"text/template"

	cronn "github.com/robfig/cron/v3"

	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/cron"
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/admin/router"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	customerRepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	walletRepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	customerService "github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go"
	"github.com/fathimasithara01/tradeverse/pkg/seeder"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	port   string
}

func InitializeApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	DB, err := migrations.ConnectDB(*cfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting db: %w", err)
	}

	seeder.CreateAdminSeeder(DB, *cfg)

	userRepo := repository.NewUserRepository(DB)
	roleRepo := repository.NewRoleRepository(DB)
	dashboardRepo := repository.NewDashboardRepository(DB)
	permissionRepo := repository.NewPermissionRepository(DB)
	activityRepo := repository.NewActivityRepository(DB)
	copyRepo := repository.NewCopyRepository(DB)
	subscriptionPlanRepo := repository.NewSubscriptionPlanRepository(DB)
	subscriptionRepo := repository.NewSubscriptionRepository(DB)
	adminWalletRepo := repository.NewAdminWalletRepository(DB)
	signalRepo := repository.NewSignalRepository(DB)

	userService := service.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	roleService := service.NewRoleService(roleRepo, permissionRepo, userRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	activityService := service.NewActivityService(activityRepo)
	copyService := service.NewCopyService(copyRepo)
	// liveSignalService := service.NewLiveSignalService(userRepo)
	subscriptionPlanService := service.NewSubscriptionPlanService(subscriptionPlanRepo)
	adminWalletService := service.NewAdminWalletService(adminWalletRepo, DB)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, subscriptionPlanRepo, userRepo, adminWalletService, DB)
	liveSignalService := service.NewLiveSignalService(signalRepo)

	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	permissionController := controllers.NewPermissionController(permissionService, roleService)
	activityController := controllers.NewActivityController(activityService)
	copyController := controllers.NewCopyController(copyService)
	// signalController := controllers.NewSignalController(liveSignalService)
	adminWalletController := controllers.NewAdminWalletController(adminWalletService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService, subscriptionPlanService)
	signalController := controllers.NewSignalController(liveSignalService)

	customerTraderSubscriptionRepo := customerRepo.NewIAdminSubscriptionRepository(DB)
	customerWalletRepo := walletRepo.NewWalletRepository(DB)
	paymentClient := paymentgateway.NewSimulatedPaymentClient()

	customerWalletService := customerService.NewWalletService(customerWalletRepo, paymentClient, DB)

	transactionRepo := repository.NewTransactionRepository(DB)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionController := controllers.NewTransactionController(transactionService)

	customerServiceForTraderSubs := customerService.NewCustomerService(
		customerTraderSubscriptionRepo,
		customerWalletService,
		customerWalletRepo,
		DB,
	)

	r := gin.Default()

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	projectRoot := filepath.Join(currentDir, "..", "..", "..")
	templatesPath := filepath.Join(projectRoot, "templates", "*.html")

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
	})

	r.LoadHTMLGlob(templatesPath)

	router.WirePublicRoutes(r, authController, signalController)
	router.WireFollowerRoutes(r, copyController, cfg)
	router.WireAdminRoutes(
		r,
		cfg,
		authController,
		dashboardController,
		userController,
		roleController,
		permissionController,
		activityController,
		roleService,
		adminWalletController,
		subscriptionController,
		transactionController,
		DB,
		signalController,
	)
	c := cronn.New()

	cron.StartCronJobs(subscriptionService, customerServiceForTraderSubs, liveSignalService, DB)
	c.AddFunc("@every 5m", func() {
		log.Println("Starting market data fetch...")
		cron.FetchAndSaveMarketData(DB)
	})
	return &App{
		engine: r,
		port:   cfg.AdminPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
