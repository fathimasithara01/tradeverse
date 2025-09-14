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
	subscriptionController *controllers.SubscriptionController,
	adminTraderSubscriptionPlanController *controllers.AdminTraderSubscriptionPlanController,

	// traderSubscriptionController *controllers.TraderSubscriptionController,

) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/signup", authController.Signup)
		public.POST("/login", authController.Login)
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

	// Subscription Routes
	subscriptionsGroup := protected.Group("/subscriptions")
	{
		subscriptionsGroup.POST("", subscriptionController.SubscribeToTrader)
		subscriptionsGroup.GET("", subscriptionController.ListMySubscriptions)
		subscriptionsGroup.GET("/:id", subscriptionController.GetSubscriptionDetails)
		subscriptionsGroup.PUT("/:id", subscriptionController.UpdateSubscription)
		subscriptionsGroup.POST("/:id/pause", subscriptionController.PauseCopyTrading)
		subscriptionsGroup.POST("/:id/resume", subscriptionController.ResumeCopyTrading)
		subscriptionsGroup.DELETE("/:id", subscriptionController.CancelSubscription)
		subscriptionsGroup.POST("/:id/simulate", subscriptionController.RunSimulation)
	}

	traderPlansGroup := protected.Group("/trader-subscription-plans")
	{
		traderPlansGroup.POST("", adminTraderSubscriptionPlanController.CreateTraderSubscriptionPlan)
		traderPlansGroup.GET("", adminTraderSubscriptionPlanController.ListTraderSubscriptionPlans)
		traderPlansGroup.GET("/:id", adminTraderSubscriptionPlanController.GetTraderSubscriptionPlanByID)
		traderPlansGroup.PUT("/:id", adminTraderSubscriptionPlanController.UpdateTraderSubscriptionPlan)
		traderPlansGroup.DELETE("/:id", adminTraderSubscriptionPlanController.DeleteTraderSubscriptionPlan)
		traderPlansGroup.PATCH("/:id/status", adminTraderSubscriptionPlanController.ToggleTraderSubscriptionPlanStatus)
	}

	// traderSubscriptionsGroup := protected.Group("")
	// {
	// 	traderSubscriptionsGroup.POST("/upgrade-to-trader", traderSubscriptionController.UpgradeToTrader)
	// 	traderSubscriptionsGroup.GET("/my-trader-subscription", traderSubscriptionController.GetMyTraderSubscription)
	// 	// Add other trader subscription management routes here if needed (e.g., renew, cancel)
	// }

	return r
}
