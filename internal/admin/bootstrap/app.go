package bootstrap

import (
	"fmt"
	"log"

	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/admin/router"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
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

	userService := service.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	roleService := service.NewRoleService(roleRepo, permissionRepo, userRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	activityService := service.NewActivityService(activityRepo)
	copyService := service.NewCopyService(copyRepo)
	liveSignalService := service.NewLiveSignalService(userRepo)
	subscriptionPlanService := service.NewSubscriptionPlanService(subscriptionPlanRepo)
	adminWalletService := service.NewAdminWalletService(adminWalletRepo, DB)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, subscriptionPlanRepo, userRepo, adminWalletService, DB)

	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	permissionController := controllers.NewPermissionController(permissionService, roleService)
	activityController := controllers.NewActivityController(activityService)
	copyController := controllers.NewCopyController(copyService)
	signalController := controllers.NewSignalController(liveSignalService)
	adminWalletController := controllers.NewAdminWalletController(adminWalletService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService, subscriptionPlanService)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

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
		signalController,
		adminWalletController,
		subscriptionController,
	)

	return &App{
		engine: r,
		port:   cfg.AdminPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
