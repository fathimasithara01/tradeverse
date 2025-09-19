package controllers

import (
	"errors" // Import errors package
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service" // Adjust based on your actual module path
	"github.com/fathimasithara01/tradeverse/pkg/models"              // Adjust based on your actual module path

	"github.com/gin-gonic/gin"
)

type TradeController struct {
	tradeService service.TradeService
}

// NewTradeController creates a new instance of TradeController
func NewTradeController(tradeService service.TradeService) *TradeController {
	return &TradeController{tradeService: tradeService}
}

// GetTraderIDFromContext is a placeholder for actual authentication logic.
// In a real application, this would extract the authenticated user's ID from the JWT token
// or session.
func GetTraderIDFromContext(c *gin.Context) (uint, error) {
	// Example: retrieve from gin.Context. You might set this in an AuthMiddleware.
	// if traderID, exists := c.Get("traderID"); exists {
	// 	if id, ok := traderID.(uint); ok {
	// 		return id, nil
	// 	}
	// }
	// return 0, errors.New("trader ID not found in context or invalid type")

	// For demonstration, returning a hardcoded ID. REPLACE WITH ACTUAL AUTH.
	return 1, nil
}

// CreateTrade handles POST /api/v1/trader/trades
// Adds a new trade (buy/sell) for the authenticated trader.
func (ctrl *TradeController) CreateTrade(c *gin.Context) {
	traderID, err := GetTraderIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input models.TradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request payload: %s", err.Error())})
		return
	}

	trade, err := ctrl.tradeService.CreateTrade(traderID, &input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTradeInput) || errors.Is(err, service.ErrInsufficientFunds) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create trade: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, trade)
}

// GetTradeByID handles GET /api/v1/trader/trades/{id}
// Retrieves details for a specific trade belonging to the authenticated trader.
func (ctrl *TradeController) GetTradeByID(c *gin.Context) {
	traderID, err := GetTraderIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade ID format"})
		return
	}

	trade, err := ctrl.tradeService.GetTradeByID(uint(tradeID), traderID)
	if err != nil {
		if errors.Is(err, service.ErrTradeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to retrieve trade: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, trade)
}

// ListTrades handles GET /api/v1/trader/trades
// Lists all trades (history + current) for the authenticated trader with pagination.
func (ctrl *TradeController) ListTrades(c *gin.Context) {
	traderID, err := GetTraderIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid pagination parameters: %s", err.Error())})
		return
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	tradesResponse, err := ctrl.tradeService.ListTrades(traderID, &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list trades: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tradesResponse)
}

// UpdateTrade handles PUT /api/v1/trader/trades/{id}
// Updates an existing trade (e.g., stop loss, take profit, or closes/cancels it) for the authenticated trader.
func (ctrl *TradeController) UpdateTrade(c *gin.Context) {
	traderID, err := GetTraderIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade ID format"})
		return
	}

	var input models.TradeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request payload: %s", err.Error())})
		return
	}

	updatedTrade, err := ctrl.tradeService.UpdateTrade(uint(tradeID), traderID, &input)
	if err != nil {
		if errors.Is(err, service.ErrTradeNotFound) || errors.Is(err, service.ErrTradeNotOpen) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrInvalidTradeInput) || errors.Is(err, service.ErrClosingPrice) || errors.New("close price is required for closing a trade").Error() == err.Error() { // Check for specific close price error
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update trade: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, updatedTrade)
}

// DeleteTrade handles DELETE /api/v1/trader/trades/{id}
// Deletes (soft delete) a trade belonging to the authenticated trader.
func (ctrl *TradeController) DeleteTrade(c *gin.Context) {
	traderID, err := GetTraderIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade ID format"})
		return
	}

	err = ctrl.tradeService.DeleteTrade(uint(tradeID), traderID)
	if err != nil {
		if errors.Is(err, service.ErrTradeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// Assuming specific error for not deletable trades
		if errors.New("open or closed trades cannot be deleted, only archived").Error() == err.Error() {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // Use 409 Conflict for business logic restriction
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete trade: %s", err.Error())})
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204 No Content for successful deletion
}
