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
	customerTraderSubController controllers.CustomerController,

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

	protected.GET("/trader-plans", customerTraderSubController.ListTraderSubscriptionPlans)
	protected.POST("/trader-plans/:plan_id/subscribe", customerTraderSubController.SubscribeToTraderPlan)
	protected.GET("/trader-subscription", customerTraderSubController.GetCustomerTraderSubscription)
	protected.POST("/trader-subscription/:subscription_id/cancel", customerTraderSubController.CancelCustomerTraderSubscription)

	return r
}
