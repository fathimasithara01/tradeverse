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
	db.ConnectDatabase()

	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*.html")

	userRepo := repository.NewUserRepository(db.DB)
	roleRepo := repository.NewRoleRepository(db.DB)
	dashboardRepo := repository.NewDashboardRepository(db.DB)

	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)

	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	dashboardController := controllers.NewDashboardController(dashboardService)

	routes.WirePublicRoutes(r, authController)

	routes.WireAdminRoutes(r, authController, dashboardController, userController, roleController)

	port := config.AppConfig.Port
	log.Printf("Server starting on port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
