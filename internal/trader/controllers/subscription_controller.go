package controllers

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/fathimasithara01/tradeverse/internal/trader/service"
// 	"github.com/gin-gonic/gin"
// )

// type SubscriptionController struct {
// 	subscriptionService service.SubscriptionService
// }

// func NewSubscriptionController(subscriptionService service.SubscriptionService) *SubscriptionController {
// 	return &SubscriptionController{subscriptionService: subscriptionService}
// }

// func (ctrl *SubscriptionController) ListTraderSubscribers(c *gin.Context) {

// 	traderIDStr := c.Query("trader_id") // Example: trader_id=1
// 	if traderIDStr == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "trader_id is required"})
// 		return
// 	}
// 	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID format"})
// 		return
// 	}

// 	subscriptions, err := ctrl.subscriptionService.ListSubscribers(uint(traderID))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, subscriptions)
// }

// func (ctrl *SubscriptionController) GetTraderSubscriberDetails(c *gin.Context) {

// 	traderIDStr := c.Query("trader_id") // Example: trader_id=1
// 	if traderIDStr == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "trader_id is required"})
// 		return
// 	}
// 	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID format"})
// 		return
// 	}

// 	subscriberID, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscriber ID format"})
// 		return
// 	}

// 	subscription, err := ctrl.subscriptionService.GetSubscriberDetails(uint(traderID), uint(subscriberID))
// 	if err != nil {
// 		if err.Error() == "subscriber not found or not associated with this trader" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, subscription)
// }
