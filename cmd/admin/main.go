package main

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	routes "github.com/fathimasithara01/tradeverse/internal/admin/router"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/seeder"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	DB, err := migrations.ConnectDB(*cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	seeder.CreateAdminSeeder(DB, *cfg)

	r := gin.Default()
	// r.LoadHTMLGlob("templates/**/*html")

	r.LoadHTMLGlob("templates/*.html")
	// r.Static("/static", "./static") // Serve static files (ensure you have a /static folder)

	userRepo := repository.NewUserRepository(DB)
	roleRepo := repository.NewRoleRepository(DB)
	dashboardRepo := repository.NewDashboardRepository(DB)
	permissionRepo := repository.NewPermissionRepository(DB)
	activityRepo := repository.NewActivityRepository(DB)
	copyRepo := repository.NewCopyRepository(DB)
	subscriptionPlanRepo := repository.NewSubscriptionPlanRepository(DB) // New
	subscriptionRepo := repository.NewSubscriptionRepository(DB)

	userService := service.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	roleService := service.NewRoleService(roleRepo, permissionRepo, userRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	activityService := service.NewActivityService(activityRepo)
	copyService := service.NewCopyService(copyRepo)
	liveSignalService := service.NewLiveSignalService(userRepo)
	subscriptionPlanService := service.NewSubscriptionPlanService(subscriptionPlanRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, subscriptionPlanRepo, userRepo)

	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	permissionController := controllers.NewPermissionController(permissionService, roleService)
	activityController := controllers.NewActivityController(activityService)
	copyController := controllers.NewCopyController(copyService)
	signalController := controllers.NewSignalController(liveSignalService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService, subscriptionPlanService)

	routes.WirePublicRoutes(r, authController, signalController)
	routes.WireFollowerRoutes(r, copyController, cfg)
	routes.WireAdminRoutes(
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
		subscriptionController,
	)

	port := cfg.AdminPort
	log.Printf("Server starting on port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
