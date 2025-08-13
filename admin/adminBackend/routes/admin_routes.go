package routes

import (
	"github.com/fathimasithara01/tradeverse/controllers"
	"github.com/fathimasithara01/tradeverse/middleware"
	"github.com/gin-gonic/gin"
)

func WireAdminRoutes(
	r *gin.Engine,
	authCtrl *controllers.AuthController,
	dashCtrl *controllers.DashboardController,
	userCtrl *controllers.UserController,
	roleCtrl *controllers.RoleController,
) {
	admin := r.Group("/admin")
	{
		// Public route for creating the first admin
		admin.GET("/register", authCtrl.ShowAdminRegisterPage)
		admin.POST("/register", authCtrl.RegisterAdmin)

		// Protected routes that require admin login
		protected := admin.Group("")
		protected.Use(middleware.JWTMiddleware())
		{
			// Dashboard Routes -> DashboardController
			protected.GET("/dashboard", dashCtrl.ShowDashboardPage)
			protected.GET("/dashboard/stats", dashCtrl.GetDashboardStats)

			// User Management Routes -> UserController
			protected.GET("/users", userCtrl.ShowUsersPage)
			protected.GET("/users/add", userCtrl.ShowAddUserPage)
			protected.GET("/users/edit/:id", userCtrl.ShowEditUserPage)
			protected.POST("/users/add", userCtrl.CreateCustomer)
			protected.POST("/users/edit/:id", userCtrl.UpdateUser)
			protected.GET("/api/users", userCtrl.GetUsers)
			protected.DELETE("/api/users/:id", userCtrl.DeleteUser)

			// Role Management Routes -> RoleController
			protected.GET("/roles", roleCtrl.ShowRolesPage)
			protected.GET("/roles/add", roleCtrl.ShowAddRolePage)
			protected.GET("/roles/edit/:id", roleCtrl.ShowEditRolePage)
			protected.POST("/roles/add", roleCtrl.CreateRole)
			protected.POST("/roles/edit/:id", roleCtrl.UpdateRole)
			protected.GET("/api/roles", roleCtrl.GetRoles)
			protected.DELETE("/api/roles/:id", roleCtrl.DeleteRole)
		}
	}
}
