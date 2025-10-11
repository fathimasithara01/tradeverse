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
	// tradeController *controllers.TradeController,
	walletCntrl *controllers.WalletController,
	// subscriberController *controllers.SubscriptionController,
	subscriberController *controllers.SubscriberController,
	liveCtrl *controllers.LiveTradeController,
	tradeSignlCntrl *controllers.SignalController,
	marketDataCnttl *controllers.MarketDataHandler,
	subsController *controllers.TraderSubscriptionController,
) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/login", authController.Login)

	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{

		protected.POST("/market-", marketDataCnttl.CreateMarketData)

		protected.GET("/trader/profile", profileController.GetTraderProfile)
		protected.POST("/trader/profile", profileController.CreateTraderProfile)
		protected.PUT("/trader/profile", profileController.UpdateTraderProfile)
		protected.DELETE("/trader/profile", profileController.DeleteTraderProfile)

		// protected.GET("/trades", tradeController.GetTraderTrades)
		// protected.POST("/trade", tradeController.CreateTrade)
		// protected.GET("/trades/:id", tradeController.GetTradeByID)
		// protected.PUT("/trade/:id", tradeController.UpdateTradeStatus)
		// protected.DELETE("/trade/:id", tradeController.RemoveTrade)

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

		protected.POST("/signals", tradeSignlCntrl.CreateSignal)
		protected.GET("/signals", tradeSignlCntrl.GetAllSignals) // Or GetSignalsByTraderID
		protected.GET("/signals/:id", tradeSignlCntrl.GetSignalByID)
		protected.PUT("/signals/:id", tradeSignlCntrl.UpdateSignal)
		protected.DELETE("/signals/:id", tradeSignlCntrl.DeleteSignal)

		protected.POST("/plans", subsController.CreateTraderSubscriptionPlan)
		protected.GET("/plans", subsController.GetMyTraderSubscriptionPlans)
		protected.GET("/plans/:planId", subsController.GetTraderSubscriptionPlanByID)
		protected.PUT("/plans/:planId", subsController.UpdateTraderSubscriptionPlan)
		protected.DELETE("/plans/:planId", subsController.DeleteTraderSubscriptionPlan)

		// protected.GET("/admin-plans", subsController.GetAllAdminSubscriptionPlans)
		// protected.POST("/subscribe-to-admin-plan/:planId", traderSubsController.SubscribeToAdminPlan)

	}

	return r
}
