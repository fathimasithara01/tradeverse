package router

import (
	"github.com/fathimasithara01/tradeverse/internal/customer/middleware" // Make sure this path is correct
	"github.com/fathimasithara01/tradeverse/internal/trader/controllers"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	tradeController *controllers.TradeController,
) *gin.Engine {
	r := gin.Default()

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{

		protected.POST("/trader/trades", tradeController.CreateTrade)
		protected.GET("/trader/trades", tradeController.ListTrades)
		protected.GET("/trader/trades/:id", tradeController.GetTradeByID)
		protected.PUT("/trader/trades/:id", tradeController.UpdateTrade)
		protected.DELETE("/trader/trades/:id", tradeController.DeleteTrade)
	}

	return r
}
