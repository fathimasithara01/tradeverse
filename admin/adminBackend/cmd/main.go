// File: adminBackend/main.go

package main

import (
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/controllers"
	"github.com/fathimasithara01/tradeverse/db"
	"github.com/fathimasithara01/tradeverse/repository"
	"github.com/fathimasithara01/tradeverse/routes"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db.ConnectDB()
	db.CreateAdminSeeder(db.DB)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	userRepo := repository.NewUserRepository(db.DB)
	roleRepo := repository.NewRoleRepository(db.DB)
	dashboardRepo := repository.NewDashboardRepository(db.DB)
	permissionRepo := repository.NewPermissionRepository(db.DB)
	activityRepo := repository.NewActivityRepository(db.DB)
	copyRepo := repository.NewCopyRepository(db.DB)

	userService := service.NewUserService(userRepo, roleRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo, userRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	activityService := service.NewActivityService(activityRepo)
	copyService := service.NewCopyService(copyRepo)
	liveSignalService := service.NewLiveSignalService(userRepo)

	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	permissionController := controllers.NewPermissionController(permissionService, roleService)
	activityController := controllers.NewActivityController(activityService)
	copyController := controllers.NewCopyController(copyService)
	signalController := controllers.NewSignalController(liveSignalService)

	routes.WirePublicRoutes(r, authController, signalController)
	routes.WireFollowerRoutes(r, copyController)
	routes.WireAdminRoutes(r, authController, dashboardController, userController, roleController, permissionController, activityController, roleService, signalController)

	port := config.AppConfig.Port
	log.Printf("Server starting on port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
