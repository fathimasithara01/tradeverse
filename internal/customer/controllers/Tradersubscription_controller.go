// internal/customer/controllers/trader_subscription.go - **NEW FILE**
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type CustomerTraderSignalSubscriptionController struct {
	subsService service.ICustomerTraderSignalSubscriptionService
}

func NewCustomerTraderSignalSubscriptionController(subsService service.ICustomerTraderSignalSubscriptionService) *CustomerTraderSignalSubscriptionController {
	return &CustomerTraderSignalSubscriptionController{subsService: subsService}
}
func (ctrl *CustomerTraderSignalSubscriptionController) GetAvailableTradersWithPlans(c *gin.Context) {
	traders, err := ctrl.subsService.GetAvailableTradersWithPlans(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch available traders: %v", err)})
		return
	}
	c.JSON(http.StatusOK, traders)
}

func (ctrl *CustomerTraderSignalSubscriptionController) SubscribeToTrader(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}
	customerID := userID.(uint)

	var input models.SubscribeToTraderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.subsService.SubscribeToTrader(c, customerID, input)
	if err != nil {
		if err.Error() == "insufficient funds in wallet" || err.Error() == "you are already subscribed to this plan" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to subscribe to trader: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully subscribed to trader's plan"})
}

func (ctrl *CustomerTraderSignalSubscriptionController) GetSignalsFromSubscribedTraders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}
	customerID := userID.(uint)

	signals, err := ctrl.subsService.GetSubscribedTradersSignals(c, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch signals from subscribed traders: %v", err)})
		return
	}

	c.JSON(http.StatusOK, signals)
}

func (ctrl *CustomerTraderSignalSubscriptionController) GetMyActiveTraderSubscriptions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}
	customerID := userID.(uint)

	subscriptions, err := ctrl.subsService.GetActiveSubscriptions(c, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch active subscriptions: %v", err)})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (ctrl *CustomerTraderSignalSubscriptionController) IsSubscribedToTrader(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}
	customerID := userID.(uint)

	traderIDParam := c.Param("traderId")
	traderID, err := strconv.ParseUint(traderIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trader ID"})
		return
	}

	isSubscribed, err := ctrl.subsService.IsCustomerSubscribedToTrader(c, customerID, uint(traderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to check subscription status: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_subscribed": isSubscribed, "trader_id": traderID})
}
