package router

import (
	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/middleware"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authController *controllers.AuthController,
	profileController *controllers.ProfileController,
	kycController *controllers.KYCController,
	walletCtrl *controllers.WalletController,
	customerTraderSubController *controllers.CustomerController,
	traderController *controllers.TraderController,
	subCtrl *controllers.SubscriptionController,
	traderWalletCtrl *controllers.TraderWalletController,
) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/signup", authController.Signup)
		public.POST("/login", authController.Login)

		public.GET("/traders", traderController.ListTraders)
		public.GET("/traders/:id", traderController.GetTraderDetails)
		public.GET("/traders/:id/performance", traderController.GetTraderPerformance)

	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/profile", profileController.GetProfile)
		protected.PUT("/profile", profileController.UpdateProfile)
		protected.DELETE("/account", profileController.DeleteAccount)
	}

	kycGroup := protected.Group("/customers")
	{
		kycGroup.POST("/kyc", kycController.SubmitKYCDocuments)
		kycGroup.GET("/kyc/status", kycController.GetKYCStatus)
	}

	walletRoutes := protected.Group("/wallet")
	{
		walletRoutes.GET("/summary", walletCtrl.GetWalletSummary)
		walletRoutes.POST("/deposit/initiate", walletCtrl.InitiateDeposit)
		walletRoutes.POST("/deposit/:deposit_id/verify", walletCtrl.VerifyDeposit)
		walletRoutes.POST("/withdraw/request", walletCtrl.RequestWithdrawal)
		walletRoutes.GET("/transactions", walletCtrl.GetWalletTransactions)
	}

	protected.GET("/subscriptions/plans", customerTraderSubController.ListTraderSubscriptionPlans)
	protected.POST("/trader-plans/:plan_id/subscribe", customerTraderSubController.SubscribeToTraderPlan)
	protected.GET("/trader-subscription", customerTraderSubController.GetCustomerTraderSubscription)
	protected.POST("/trader-subscription/:subscription_id/cancel", customerTraderSubController.CancelCustomerTraderSubscription)

	protected.POST("/subscribe", traderWalletCtrl.SubscribeCustomer)
	protected.GET("/balance", traderWalletCtrl.GetBalance)

	return r
}
