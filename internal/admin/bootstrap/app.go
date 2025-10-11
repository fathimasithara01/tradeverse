package bootstrap

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"text/template"

	cronn "github.com/robfig/cron/v3"

	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/cron"
	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/admin/router"
	adminService "github.com/fathimasithara01/tradeverse/internal/admin/service"

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

	adminUserRepo := adminRepo.NewUserRepository(DB)
	adminRoleRepo := adminRepo.NewRoleRepository(DB)
	adminDashboardRepo := adminRepo.NewDashboardRepository(DB)
	adminPermissionRepo := adminRepo.NewPermissionRepository(DB)
	adminActivityRepo := adminRepo.NewActivityRepository(DB)
	adminSubscriptionPlanRepo := adminRepo.NewSubscriptionPlanRepository(DB)
	adminSubscriptionRepo := adminRepo.NewSubscriptionRepository(DB)
	adminAdminWalletRepo := adminRepo.NewAdminWalletRepository(DB)
	adminSignalRepo := adminRepo.NewSignalRepository(DB)
	adminTransactionRepo := adminRepo.NewTransactionRepository(DB)

	adminUserService := adminService.NewUserService(adminUserRepo, adminRoleRepo, cfg.JWTSecret)
	adminRoleService := adminService.NewRoleService(adminRoleRepo, adminPermissionRepo, adminUserRepo)
	adminDashboardService := adminService.NewDashboardService(adminDashboardRepo)
	adminPermissionService := adminService.NewPermissionService(adminPermissionRepo)
	adminActivityService := adminService.NewActivityService(adminActivityRepo)
	adminSubscriptionPlanService := adminService.NewSubscriptionPlanService(adminSubscriptionPlanRepo)
	adminAdminWalletService := adminService.NewAdminWalletService(adminAdminWalletRepo, DB)
	adminSubscriptionService := adminService.NewSubscriptionService(adminSubscriptionRepo, adminSubscriptionPlanRepo, adminUserRepo, adminAdminWalletService, DB)
	adminLiveSignalService := adminService.NewLiveSignalService(adminSignalRepo)
	adminTransactionService := adminService.NewTransactionService(adminTransactionRepo)

	adminAuthController := controllers.NewAuthController(adminUserService)
	adminUserController := controllers.NewUserController(adminUserService)
	adminRoleController := controllers.NewRoleController(adminRoleService)
	adminDashboardController := controllers.NewDashboardController(adminDashboardService)
	adminPermissionController := controllers.NewPermissionController(adminPermissionService, adminRoleService)
	adminActivityController := controllers.NewActivityController(adminActivityService)
	adminAdminWalletController := controllers.NewAdminWalletController(adminAdminWalletService)
	adminSubscriptionController := controllers.NewSubscriptionController(adminSubscriptionService, adminSubscriptionPlanService)
	adminSignalController := controllers.NewSignalController(adminLiveSignalService)
	adminTransactionController := controllers.NewTransactionController(adminTransactionService)

	customerAdminSubscriptionRepository := customerRepo.NewIAdminSubscriptionRepository(DB)
	customerWalletConcreteRepo := walletRepo.NewWalletRepository(DB)
	customerSubscriptionPlanRepo := customerRepo.NewCustomerTraderSubscriptionRepository(DB)

	paymentClient := paymentgateway.NewSimulatedPaymentClient()

	customerWalletService := customerService.NewWalletService(DB, customerWalletConcreteRepo, paymentClient)
	_ = customerService.NewCustomerTraderSubscriptionService(
		customerSubscriptionPlanRepo,
		DB,
	)

	customerAdminSubscriptionServiceForCron := customerService.NewAdminSubscriptionService(
		customerAdminSubscriptionRepository,
		customerWalletService,
		customerWalletConcreteRepo,
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

	router.WireAdminRoutes(
		r,
		cfg,
		adminAuthController,
		adminDashboardController,
		adminUserController,
		adminRoleController,
		adminPermissionController,
		adminActivityController,
		adminRoleService,
		adminAdminWalletController,
		adminSubscriptionController,
		adminTransactionController,
		DB,
		adminSignalController,
	)

	c := cronn.New()
	cron.StartCronJobs(adminSubscriptionService, customerAdminSubscriptionServiceForCron, adminLiveSignalService, DB)
	c.AddFunc("@every 5m", func() {
		log.Println("Starting market data fetch...")
		cron.FetchAndSaveMarketData(DB)
	})
	c.Start()

	return &App{
		engine: r,
		port:   cfg.AdminPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
