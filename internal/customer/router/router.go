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
	walletController *controllers.WalletController,
	subscriptionController *controllers.CustomerSubscriptionController,

) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/signup", authController.Signup)
		public.POST("/login", authController.Login)

		public.GET("/subscriptions/plans", subscriptionController.GetAllTraderSubscriptionPlans)
		public.GET("/subscriptions/plans/:id", subscriptionController.GetSubscriptionPlanDetails)
		// Simulation can also be public if you allow anonymous simulation
		public.POST("/subscriptions/plans/:plan_id/simulate", subscriptionController.SimulateSubscription)
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

	protected.GET("/wallet", walletController.GetWalletSummary)
	protected.POST("/wallet/deposit", walletController.CreateDepositRequest)
	protected.POST("/wallet/deposit/verify", walletController.VerifyDeposit)
	protected.POST("/wallet/withdraw", walletController.CreateWithdrawalRequest)
	protected.GET("/wallet/transactions", walletController.ListWalletTransactions)

	// Subscribe to trader
	protected.POST("/subscriptions", subscriptionController.SubscribeToTrader)

	// Get my subscriptions
	protected.GET("/subscriptions", subscriptionController.ListMySubscriptions)
	protected.GET("/subscriptions/:id", subscriptionController.GetMySubscriptionDetails)

	// Update subscription (allocation/risk)
	protected.PUT("/subscriptions/:id", subscriptionController.UpdateSubscriptionSettings)

	// Pause/Resume copy trading
	protected.POST("/subscriptions/:id/pause", subscriptionController.PauseCopyTrading)
	protected.POST("/subscriptions/:id/resume", subscriptionController.ResumeCopyTrading)

	// Cancel subscription
	protected.DELETE("/subscriptions/:id", subscriptionController.CancelSubscription)

	return r
}
