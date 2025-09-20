package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type TradeController struct {
	svc service.TradeService
}

func NewTradeController(svc service.TradeService) *TradeController {
	return &TradeController{svc: svc}
}

func (t *TradeController) ListTrades(c *gin.Context) {
	traderID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := t.svc.ListTrades(traderID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trades"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (t *TradeController) GetTrade(c *gin.Context) {
	traderID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	trade, err := t.svc.GetTrade(uint(id), traderID)
	if err != nil || trade == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trade not found"})
		return
	}
	c.JSON(http.StatusOK, trade)
}

func (c *TradeController) CreateTrade(ctx *gin.Context) {
	var req models.TradeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	traderID, exists := ctx.Get("trader_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req.TraderID = traderID.(uint) // force override

	trade, err := c.svc.CreateTrade(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, trade)
}

func (t *TradeController) UpdateTrade(c *gin.Context) {
	traderID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var input models.TradeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade, err := t.svc.UpdateTrade(uint(id), traderID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trade)
}

func (t *TradeController) DeleteTrade(c *gin.Context) {
	traderID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	if err := t.svc.DeleteTrade(uint(id), traderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete trade"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trade deleted"})
}
