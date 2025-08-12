package routes

import (
	"github.com/fathimasithara01/tradeverse/controllers"
	"github.com/fathimasithara01/tradeverse/db"
	"github.com/fathimasithara01/tradeverse/middleware"
	"github.com/fathimasithara01/tradeverse/repository"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)

	dashboardRepo := repository.NewDashboardRepository(db.DB)
	dashboardService := service.NewDashboardService(dashboardRepo)

	userController := controllers.NewUserController(userService, dashboardService)
	// --- ROUTE DEFINITIONS ---
	admin := r.Group("/admin")
	{
		admin.GET("/register", userController.ShowRegisterPage)
		admin.GET("/login", userController.ShowLoginPage)
		admin.POST("/login", userController.LoginAdmin)

		protected := admin.Group("")
		protected.Use(middleware.JWTMiddleware())
		{
			protected.GET("/dashboard", userController.ShowDashboardPage)
			protected.GET("/dashboard/stats", userController.GetDashboardStats)
			// protected.GET("/api/monthly-orders", adminController.GetMonthlyOrderStats)

			// protected.GET("/users", userController.ShowUsersPage)
			// protected.POST("/users/add", userController.CreateCustomer)
			// protected.GET("/users/add", userController.ShowAddUserPage)

			// protected.POST("/users/edit/:id", customerController.UpdateCustomer) // Handle form submission
			// protected.DELETE("/api/users/:id", customerController.DeleteCustomer)
		}
	}
}
