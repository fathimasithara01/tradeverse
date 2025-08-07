package routes

import (
	"github.com/fathimasithara01/tradeverse/admin/controllers"
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		ctrl := controllers.NewAdminController(db.DB)

		r.GET("/admin/register", ctrl.ShowRegisterPage)
		r.POST("/admin/register", ctrl.RegisterAdmin)

		r.GET("/admin/login", ctrl.ShowLoginPage)
		r.POST("/admin/login", ctrl.LoginAdmin)

		r.GET("/admin/dashboard", ctrl.AdminDashboard)
		r.GET("/admin/logout", ctrl.LogoutAdmin)
		// admin.POST("/register", controllers.RegisterAdmin)
		// admin.POST("/login", controllers.LoginAdmin)

		// admin.GET("/dashboard", middleware.JWTMiddleware(), func(c *gin.Context) {
		// 	c.JSON(200, gin.H{"message": "Welcome to Admin Dashboard"})
		// })

		admin.GET("/dashboard", middleware.JWTMiddleware(), controllers.GetAdminDashboard)
		// secured.GET("/dashboard", func(c *gin.Context) {
		//     c.JSON(200, gin.H{"message": "Welcome to Admin Dashboard!"})
		// })

		admin.GET("/users", middleware.JWTMiddleware(), controllers.GetAllUsers)
		admin.GET("/traders", middleware.JWTMiddleware(), controllers.GetAllTraders)
		admin.PUT("/trader/:id/ban", middleware.JWTMiddleware(), controllers.ToggleBanTrader)

		admin.GET("/plans/pending", middleware.JWTMiddleware(), controllers.GetPendingPlans)
		admin.PUT("/plans/:id/approve", middleware.JWTMiddleware(), controllers.ApprovePlan)
		admin.PUT("/plans/:id/reject", middleware.JWTMiddleware(), controllers.RejectPlan)

		admin.GET("/signals", middleware.JWTMiddleware(), controllers.GetAllSignals)

		admin.GET("/subscriptions", middleware.JWTMiddleware(), controllers.GetAllSubscriptions)

		admin.GET("/payments", middleware.JWTMiddleware(), controllers.GetAllPayments)

		admin.GET("/revenue/monthly", middleware.JWTMiddleware(), controllers.GetMonthlyRevenue)

		admin.PUT("/signals/:id/deactivate", middleware.JWTMiddleware(), controllers.DeactivateSignal)

		admin.POST("/announcement", middleware.JWTMiddleware(), controllers.CreateAnnouncement)
		admin.GET("/announcements", middleware.JWTMiddleware(), controllers.GetAllAnnouncements)

		admin.PUT("/users/:id/ban", middleware.JWTMiddleware(), controllers.BanUser)
		admin.PUT("/users/:id/unban", middleware.JWTMiddleware(), controllers.UnbanUser)

		admin.GET("/stats/plans", middleware.JWTMiddleware(), controllers.GetPlanStats)

		admin.GET("/logs", middleware.JWTMiddleware(), controllers.GetAllLogs)

		// admin.GET("/export/users/csv", middleware.JWTMiddleware(), controllers.ExportUsersCSV)
		// admin.GET("/export/subscriptions/pdf", middleware.JWTMiddleware(), controllers.ExportSubscriptionsPDF)

		admin.GET("/signals/pending", middleware.JWTMiddleware(), controllers.GetPendingSignals)
		admin.PUT("/signals/:id/approve", middleware.JWTMiddleware(), controllers.ApproveSignal)
		admin.PUT("/signals/:id/reject", middleware.JWTMiddleware(), controllers.RejectSignal)

		admin.GET("/plans", middleware.JWTMiddleware(), controllers.GetAllPlans)
		admin.POST("/plans", middleware.JWTMiddleware(), controllers.CreatePlan)
		admin.PUT("/plans/:id", middleware.JWTMiddleware(), controllers.UpdatePlan)
		admin.DELETE("/plans/:id", middleware.JWTMiddleware(), controllers.DeletePlan)

		admin.GET("/traders/analytics", middleware.JWTMiddleware(), controllers.GetTraderStats)
		admin.GET("/traders/rankings", middleware.JWTMiddleware(), controllers.GetTopRankedTraders)

		admin.GET("/withdrawals/pending", middleware.JWTMiddleware(), controllers.GetPendingWithdrawals)
		admin.PUT("/withdrawals/:id/approve", middleware.JWTMiddleware(), controllers.ApproveWithdrawal)
		admin.PUT("/withdrawals/:id/reject", middleware.JWTMiddleware(), controllers.RejectWithdrawal)

		admin.GET("/wallets/:user_id", middleware.JWTMiddleware(), controllers.GetWalletDetails)
		admin.POST("/wallets/:user_id/credit", middleware.JWTMiddleware(), controllers.CreditWallet)
		admin.POST("/wallets/:user_id/debit", middleware.JWTMiddleware(), controllers.DebitWallet)

		admin.POST("/notifications", middleware.JWTMiddleware(), controllers.SendNotification)

		admin.GET("/payments/history", middleware.JWTMiddleware(), controllers.AdminPaymentHistory)

		admin.GET("/analytics/signals", controllers.GetSignalAnalytics)
		admin.GET("/analytics/traders", controllers.GetTraderPerformance)

		admin.GET("/revenue/splits", controllers.GetAllRevenueSplits)

		admin.GET("/traders/ranking", controllers.GetTraderRankings)
		admin.GET("/trader/:id/badge", controllers.GetTraderBadge)

	}

}
