package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/gin-gonic/gin"
)

type MarketDataHandler struct {
	service service.MarketDataService
}

func NewMarketDataHandler(service service.MarketDataService) *MarketDataHandler {
	return &MarketDataHandler{service: service}
}

type MarketDataRequest struct {
	Symbol string  `json:"symbol" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

func (h *MarketDataHandler) CreateMarketData(c *gin.Context) {
	var req MarketDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.CreateMarketData(req.Symbol, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create market data"})
		return
	}

	c.JSON(http.StatusOK, data)
}
