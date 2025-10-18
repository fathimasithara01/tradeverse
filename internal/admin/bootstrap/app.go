// internal/admin/bootstrap/app.go
package bootstrap

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"text/template"

	// Keep this import for the explicit cron.New()

	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/cron"
	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/admin/router"
	adminService "github.com/fathimasithara01/tradeverse/internal/admin/service"

	// These customer-side imports are typically not needed in admin bootstrap unless
	// admin needs to *directly* instantiate and use customer-side services for something
	// other than the admin cron, or if an admin controller depends on them.
	// Based on the given admin code, they don't seem necessary for the core admin app.
	// Leaving them commented out or removing them entirely if unused is best practice.
	cusSvc "github.com/fathimasithara01/tradeverse/internal/customer/service"
	// walletRepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	// customerService "github.com/fathimasithara01/tradeverse/internal/customer/service"

	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"

	// paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go" // Not directly used in admin bootstrap
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

	// Admin Repositories
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

	// Admin Services
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

	// Admin Controllers
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

	var customerAdminUpgradeSubscriptionService cusSvc.CustomerSubscriptionService
	_ = customerAdminUpgradeSubscriptionService

	r := gin.Default()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get current file path")
	}
	currentDir := filepath.Dir(filename)
	projectRoot := filepath.Join(currentDir, "..", "..", "..")
	templatesPath := filepath.Join(projectRoot, "templates", "*.html")
	staticPath := filepath.Join(projectRoot, "static")

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			if a < b {
				return 0
			}
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
	r.Static("/static", staticPath)

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

	cron.StartCronJobs(adminSubscriptionService, customerAdminUpgradeSubscriptionService, adminLiveSignalService, DB)

	return &App{
		engine: r,
		port:   cfg.AdminPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
