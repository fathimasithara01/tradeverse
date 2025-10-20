package router

import (
	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/middleware"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	adminWalletController *controllers.AdminWalletController,
	subscriptionController *controllers.SubscriptionController,
	tranasactionController *controllers.TransactionController,
	db *gorm.DB,
	signalCtrl *controllers.SignalController,

) {
	authz := middleware.NewAuthzMiddleware(roleService)

	admin := r.Group("/admin")
	{
		admin.GET("/login", authCtrl.ShowLoginPage)
		admin.POST("/login", authCtrl.LoginUser)
	}

	{
		{
			admin.Use(middleware.DBMiddleware(db))

			admin.GET("/signal-cards", controllers.GetSignalCardsPage)
			admin.GET("/api/market-data", controllers.GetMarketDataAPI)

			admin.GET("/signals/create", signalCtrl.ShowCreateSignalCardPage)
			admin.POST("/api/signals", signalCtrl.CreateSignal)

			protected := admin.Group("")
			protected.Use(middleware.JWTMiddleware(cfg))
			{
				protected.GET("/dashboard", authz.RequirePermission("view_dashboard"), dashCtrl.ShowDashboardPage)
				protected.GET("/dashboard/stats", dashCtrl.GetDashboardStats)
				protected.GET("/dashboard/charts", dashCtrl.GetChartData)
				protected.GET("/dashboard/top-traders", dashCtrl.GetTopTraders)
				protected.GET("/dashboard/latest-signups", dashCtrl.GetLatestSignups)
				protected.GET("/dashboard/market-data", dashCtrl.GetLiveMarketData)

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

				protected.GET("/financials/subscriptions", subscriptionController.ShowSubscriptionsPage)
				protected.GET("/financials/api/subscriptions", subscriptionController.GetSubscriptions)
				protected.GET("/financials/subscription-plans", subscriptionController.ShowSubscriptionPlansPage)
				protected.GET("/financials/api/subscription-plans", subscriptionController.GetSubscriptionPlans)
				protected.POST("/financials/api/subscription-plans", subscriptionController.CreateSubscriptionPlan)
				protected.PUT("/financials/api/subscription-plans/:id", subscriptionController.UpdateSubscriptionPlan)
				protected.DELETE("/financials/api/subscription-plans/:id", subscriptionController.DeleteSubscriptionPlan)
				protected.GET("/financials/api/subscription-plans/:id", subscriptionController.GetSubscriptionPlanByID)
				protected.PUT("/financials/api/traders/:id/status", subscriptionController.UpdateTraderStatus)

				protected.GET("/financials/wallet", adminWalletController.ShowAdminWalletPage)
				protected.GET("/financials/api/wallet/summary", adminWalletController.GetAdminWalletSummary)
				protected.POST("/financials/api/wallet/deposit", adminWalletController.AdminInitiateDeposit)
				protected.POST("/financials/api/wallet/deposit/:deposit_id/verify", adminWalletController.AdminVerifyDeposit)
				protected.POST("/financials/api/wallet/withdraw", adminWalletController.AdminRequestWithdrawal)
				protected.GET("/financials/api/wallet/transactions", adminWalletController.AdminGetWalletTransactions)
				protected.GET("/financials/api/transactions/all", adminWalletController.AdminGetAllPlatformTransactions)

				protected.GET("/financials/api/withdrawals/pending", adminWalletController.GetPendingWithdrawals)
				protected.POST("/financials/api/withdrawals/:id/action", adminWalletController.AdminApproveOrRejectWithdrawal)

				protected.GET("/transactions", tranasactionController.GetTransactionsPage)
				protected.GET("/api/transactions", tranasactionController.GetTransactionsAPI)
			}

		}

	}
}
