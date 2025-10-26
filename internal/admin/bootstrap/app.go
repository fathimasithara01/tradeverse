package bootstrap

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/cron"
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/admin/router"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"

	cusSvc "github.com/fathimasithara01/tradeverse/internal/customer/service"

	"github.com/fathimasithara01/tradeverse/migrations"

	"github.com/fathimasithara01/tradeverse/pkg/models"
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

	if err := models.EnsureDefaultWebConfiguration(DB); err != nil {
		return nil, fmt.Errorf("error ensuring default web configuration: %w", err)
	}

	seeder.CreateAdminSeeder(DB, *cfg)

	// Admin Repositories
	adminUserRepo := repository.NewUserRepository(DB)
	adminRoleRepo := repository.NewRoleRepository(DB)
	adminDashboardRepo := repository.NewDashboardRepository(DB)
	adminPermissionRepo := repository.NewPermissionRepository(DB)
	adminActivityRepo := repository.NewActivityRepository(DB)
	adminSubscriptionPlanRepo := repository.NewSubscriptionPlanRepository(DB)
	adminSubscriptionRepo := repository.NewSubscriptionRepository(DB)
	adminAdminWalletRepo := repository.NewAdminWalletRepository(DB)
	adminSignalRepo := repository.NewSignalRepository(DB)
	adminTransactionRepo := repository.NewTransactionRepository(DB)
	commissionRepo := repository.NewCommissionRepository(DB)
	adminWebConfigRepo := repository.NewWebConfigurationRepository(DB)

	// Admin Services
	adminUserService := service.NewUserService(adminUserRepo, adminRoleRepo, cfg.JWT.Secret)
	adminRoleService := service.NewRoleService(adminRoleRepo, adminPermissionRepo, adminUserRepo)
	adminDashboardService := service.NewDashboardService(adminDashboardRepo)
	adminPermissionService := service.NewPermissionService(adminPermissionRepo)
	adminActivityService := service.NewActivityService(adminActivityRepo)
	adminSubscriptionPlanService := service.NewSubscriptionPlanService(adminSubscriptionPlanRepo)
	adminAdminWalletService := service.NewAdminWalletService(adminAdminWalletRepo, DB)
	adminSubscriptionService := service.NewSubscriptionService(adminSubscriptionRepo, adminSubscriptionPlanRepo, adminUserRepo, adminAdminWalletService, DB)
	adminLiveSignalService := service.NewLiveSignalService(adminSignalRepo)
	adminTransactionService := service.NewTransactionService(adminTransactionRepo)
	marketDataService := service.NewMarketDataService()
	commissionService := service.NewCommissionService(commissionRepo, DB)
	adminWebConfigService := service.NewWebConfigurationService(adminWebConfigRepo)

	// Admin Controllers
	adminAuthController := controllers.NewAuthController(adminUserService)
	adminUserController := controllers.NewUserController(adminUserService)
	adminRoleController := controllers.NewRoleController(adminRoleService)
	adminDashboardController := controllers.NewDashboardController(adminDashboardService, marketDataService)
	adminPermissionController := controllers.NewPermissionController(adminPermissionService, adminRoleService)
	adminActivityController := controllers.NewActivityController(adminActivityService)
	adminAdminWalletController := controllers.NewAdminWalletController(adminAdminWalletService)
	adminSubscriptionController := controllers.NewSubscriptionController(adminSubscriptionService, adminSubscriptionPlanService)
	adminSignalController := controllers.NewSignalController(adminLiveSignalService)
	adminTransactionController := controllers.NewTransactionController(adminTransactionService)
	commissionController := controllers.NewCommissionController(commissionService)
	adminWebConfigController := controllers.NewWebConfigurationController(adminWebConfigService)

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
		commissionController,
		adminWebConfigController,
	)

	cron.StartCronJobs(adminSubscriptionService, customerAdminUpgradeSubscriptionService, adminLiveSignalService, DB)

	return &App{
		engine: r,
		port:   cfg.Server.AdminPort,
	}, nil
}

func (a *App) Engine() *gin.Engine {
	return a.engine
}
func (a *App) Run() error {
	log.Printf("Server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
