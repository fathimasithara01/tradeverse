package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type LiveTradeController struct {
	svc service.LiveTradeService
}

func NewLiveTradeController(svc service.LiveTradeService) *LiveTradeController {
	return &LiveTradeController{svc: svc}
}

// POST /api/v1/trader/live
func (c *LiveTradeController) PublishLiveTrade(ctx *gin.Context) {
	var req models.LiveTrade
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	traderID, _ := ctx.Get("userID")
	req.TraderID = traderID.(uint)
	req.Status = "OPEN"

	if err := c.svc.PublishLiveTrade(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish live trade"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Live trade published", "trade": req})
}

// GET /api/v1/trader/live
func (c *LiveTradeController) GetActiveTrades(ctx *gin.Context) {
	traderID, _ := ctx.Get("userID")
	trades, err := c.svc.GetActiveTrades(traderID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active trades"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"active_trades": trades})
}
