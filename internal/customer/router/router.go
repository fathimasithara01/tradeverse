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
	traderController *controllers.TraderController,
	custmerTraderSignlsController *controllers.CustomerTraderSignalSubscriptionController,
	subscriptionPlanController *controllers.SubscriptionPlanController,
) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/signup", authController.Signup)
		public.POST("/login", authController.Login)

		public.GET("/traders", traderController.ListTraders)
		public.GET("/traders/:trader_id", traderController.GetTraderDetails)
		public.GET("/traders/:trader_id/performance", traderController.GetTraderPerformance)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{

		protected.GET("/subscription-plans", subscriptionPlanController.GetAllSubscriptionPlans)
		protected.GET("/subscription-plans/:id", subscriptionPlanController.GetSubscriptionPlanByID)
		protected.POST("/subscription-plans/:id/subscribe", subscriptionPlanController.SubscribeToPlan)
		protected.DELETE("/my-subscriptions/:id", subscriptionPlanController.CancelSubscription)
		protected.GET("/my-subscriptions", subscriptionPlanController.GetUserSubscriptions)

		protected.GET("/profile", profileController.GetProfile)
		protected.PUT("/profile", profileController.UpdateProfile)
		protected.DELETE("/account", profileController.DeleteAccount)

		protected.GET("/traders/plans", custmerTraderSignlsController.GetAvailableTradersWithPlans)
		protected.POST("/subscribe", custmerTraderSignlsController.SubscribeToTrader)
		protected.GET("/signals", custmerTraderSignlsController.GetSignalsFromSubscribedTraders)
		protected.GET("/my-trader-subscriptions", custmerTraderSignlsController.GetMyActiveTraderSubscriptions)
		protected.GET("/subscribed-to-trader/:traderId", custmerTraderSignlsController.IsSubscribedToTrader)

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
	}

	return r
}
