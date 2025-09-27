package router

import (
	"github.com/fathimasithara01/tradeverse/internal/customer/middleware"
	"github.com/fathimasithara01/tradeverse/internal/trader/controllers"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authController *controllers.AuthController,
	profileController *controllers.TraderProfileController,
	tradeController *controllers.TradeController,
	walletCntrl *controllers.WalletController,
	// subscriberController *controllers.SubscriptionController,
	subscriberController *controllers.SubscriberController,
	liveCtrl *controllers.LiveTradeController,

) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/login", authController.Login)

	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{

		protected.GET("/trader/profile", profileController.GetTraderProfile)
		protected.POST("/trader/profile", profileController.CreateTraderProfile)
		protected.PUT("/trader/profile", profileController.UpdateTraderProfile)
		protected.DELETE("/trader/profile", profileController.DeleteTraderProfile)

		protected.GET("/trades", tradeController.GetTraderTrades)
		protected.POST("/trade", tradeController.CreateTrade)
		protected.GET("/trades/:id", tradeController.GetTradeByID)
		protected.PUT("/trade/:id", tradeController.UpdateTradeStatus)
		protected.DELETE("/trade/:id", tradeController.RemoveTrade)

		protected.GET("/wallet", walletCntrl.GetBalance)
		protected.POST("/wallet/deposit", walletCntrl.Deposit)
		protected.POST("/wallet/withdraw", walletCntrl.Withdraw)
		protected.GET("/wallet/transactions", walletCntrl.TransactionHistory)

		protected.GET("/trader/subscribers", subscriberController.ListSubscribers)
		protected.GET("/trader/subscribers/:id", subscriberController.GetSubscriber)

		// protected.GET("/trader/subscribers", subscriberController.ListTraderSubscribers)
		// protected.GET("/trader/subscribers/:id", subscriberController.GetTraderSubscriberDetails)

		protected.POST("/trader/live", liveCtrl.PublishLiveTrade)
		protected.GET("/trader/live", liveCtrl.GetActiveTrades)

	}

	return r
}
