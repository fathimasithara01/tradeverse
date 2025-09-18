package routes

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/middleware"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

func WireAdminRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authCtrl *controllers.AuthController,
	dashCtrl *controllers.DashboardController,
	userCtrl *controllers.UserController,
	roleCtrl *controllers.RoleController,
	permCtrl *controllers.PermissionController,
	activityCtrl *controllers.ActivityController,
	roleService service.IRoleService,
	signalCtrl *controllers.SignalController,
	adminWalletController *controllers.AdminWalletController,
	subscriptionController *controllers.SubscriptionController,

) {
	authz := middleware.NewAuthzMiddleware(roleService)

	admin := r.Group("/admin")
	{

		protected := admin.Group("")
		protected.Use(middleware.JWTMiddleware(cfg))

		{
			protected.GET("/dashboard", authz.RequirePermission("view_dashboard"), dashCtrl.ShowDashboardPage)
			protected.GET("/dashboard/stats", dashCtrl.GetDashboardStats)
			protected.GET("/dashboard/charts", dashCtrl.GetChartData)
			protected.GET("/dashboard/top-traders", dashCtrl.GetTopTraders)
			protected.GET("/dashboard/latest-signups", dashCtrl.GetLatestSignups)

			protected.GET("/api/users/advanced", userCtrl.GetAllUsersAdvanced)

			protected.GET("/signals", signalCtrl.ShowLiveSignalsPage)
			protected.GET("/api/signals", signalCtrl.GetLiveSignals)

			protected.GET("/api/users/all", userCtrl.GetAllUsers)
			protected.GET("/users/all", authz.RequirePermission("manage_users"), userCtrl.ShowUsersPage)
			protected.GET("/users/internal/add", userCtrl.ShowAddInternalUserPage)
			protected.POST("/users/internal/add", userCtrl.CreateInternalUser)
			protected.GET("/users/edit/:id", userCtrl.ShowEditUserPage)

			protected.GET("/users/add", userCtrl.ShowAddCustomerPage)
			protected.POST("/users/add", userCtrl.CreateCustomer)
			protected.POST("/users/edit/:id", userCtrl.UpdateUser)

			protected.GET("/roles", authz.RequirePermission("manage_roles"), roleCtrl.ShowRolesPage)
			protected.GET("/roles/add", roleCtrl.ShowAddRolePage)
			protected.GET("/roles/edit/:id", roleCtrl.ShowEditRolePage)
			protected.POST("/roles/add", roleCtrl.CreateRole)
			protected.POST("/roles/edit/:id", roleCtrl.UpdateRole)
			protected.GET("/api/roles", roleCtrl.GetRoles)
			protected.DELETE("/api/roles/:id", roleCtrl.DeleteRole)

			protected.GET("/roles/permissions", permCtrl.ShowAssignPage)
			protected.GET("/api/permissions", permCtrl.GetAllPermissions)
			protected.GET("/api/roles/:id/permissions", permCtrl.GetPermissionsForRole)
			protected.POST("/api/roles/:id/permissions", permCtrl.AssignPermissionsToRole)

			protected.GET("/users/customers", userCtrl.ShowCustomersPage)
			protected.GET("/users/traders", userCtrl.ShowTradersPage)
			protected.GET("/api/users/customers", userCtrl.GetCustomers)
			protected.GET("/api/users/traders", userCtrl.GetTraders)

			protected.GET("/users/traders/add", userCtrl.ShowAddTraderPage)
			protected.POST("/users/traders/add", userCtrl.CreateTrader)

			protected.GET("/users/traders/approval", userCtrl.ShowTraderApprovalPage)

			protected.GET("/api/users/traders/pending", userCtrl.GetPendingTraders)
			protected.GET("/api/users/traders/approved", userCtrl.GetApprovedTraders)
			protected.POST("/api/users/traders/:id/approve", userCtrl.ApproveTrader)
			protected.POST("/api/users/traders/:id/reject", userCtrl.RejectTrader)

			protected.GET("/users/assign-role", authz.RequirePermission("manage_roles"), userCtrl.ShowAssignRolePage)

			protected.GET("/api/users/for-role-assignment", userCtrl.GetUsersForRoleAssignment)
			protected.POST("/api/users/assign-role", authz.RequirePermission("manage_roles"), userCtrl.AssignRoleToUser)
			protected.DELETE("/api/users/:id", authz.RequirePermission("delete_users"), userCtrl.DeleteUser)

			protected.GET("/activity/live", activityCtrl.ShowLiveCopyingPage)
			protected.GET("/activity/logs", activityCtrl.ShowTradeErrorsPage)

			protected.GET("/api/activity/live", activityCtrl.GetActiveSessions)
			protected.GET("/api/activity/logs", activityCtrl.GetTradeLogs)

			financials := admin.Group("/financials")
			{
				financials.GET("/subscriptions", subscriptionController.ShowSubscriptionsPage)
				financials.GET("/api/subscriptions", subscriptionController.GetSubscriptions)

				financials.GET("/subscription-plans", subscriptionController.ShowSubscriptionPlansPage)
				financials.GET("/api/subscription-plans", subscriptionController.GetSubscriptionPlans)
				financials.POST("/api/subscription-plans", subscriptionController.CreateSubscriptionPlan)
				financials.PUT("/api/subscription-plans/:id", subscriptionController.UpdateSubscriptionPlan)
				financials.DELETE("/api/subscription-plans/:id", subscriptionController.DeleteSubscriptionPlan)

				financials.GET("/wallet", adminWalletController.ShowAdminWalletPage)
				financials.GET("/api/wallet/summary", adminWalletController.GetAdminWalletSummary)
				financials.POST("/api/wallet/deposit", adminWalletController.AdminInitiateDeposit)
				financials.POST("/api/wallet/deposit/:deposit_id/verify", adminWalletController.AdminVerifyDeposit) // Or integrate with a real webhook
				financials.POST("/api/wallet/withdraw", adminWalletController.AdminRequestWithdrawal)
				financials.GET("/api/wallet/transactions", adminWalletController.AdminGetWalletTransactions)

				financials.GET("/api/withdrawals/pending", adminWalletController.GetPendingWithdrawals)
				financials.POST("/api/withdrawals/:id/action", adminWalletController.AdminApproveOrRejectWithdrawal)
			}
		}

	}
}
