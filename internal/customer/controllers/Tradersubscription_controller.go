package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type CustomerSubscriptionController struct {
	TraderSubscriptionService service.ITraderSubscriptionService // Renamed for clarity
}

func NewCustomerSubscriptionController(svc service.ITraderSubscriptionService) *CustomerSubscriptionController {
	return &CustomerSubscriptionController{TraderSubscriptionService: svc}
}

func (c *CustomerSubscriptionController) SubscribeTrader(ctx *gin.Context) {
	var req models.TraderSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.TraderSubscriptionService.SubscribeCustomer(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// func (ctrl *CustomerSubscriptionController) SubscribeToTrader(c *gin.Context) {
// 	customerID := c.MustGet("userID").(uint) // Assuming userID is set by auth middleware

// 	traderIDStr := c.Param("trader_id")
// 	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid trader ID format"})
// 		return
// 	}

// 	var req models.TraderSubscriptionRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid request payload: %v", err.Error())})
// 		return
// 	}
// 	req.TraderID = uint(traderID)
// 	req.CustomerID = customerID

// 	response, err := ctrl.TraderSubscriptionService.SubscribeCustomerToTrader(c.Request.Context(), req)
// 	if err != nil {
// 		statusCode := http.StatusInternalServerError
// 		if errors.Is(err, service.ErrTraderNotFound) || errors.Is(err, service.ErrSubscriptionPlanNotFound) {
// 			statusCode = http.StatusNotFound
// 		} else if errors.Is(err, service.ErrInsufficientFunds) || errors.Is(err, service.ErrAlreadySubscribed) {
// 			statusCode = http.StatusConflict // Use 409 for conflict states like already subscribed or insufficient funds
// 		} else if errors.Is(err, service.ErrCustomerWalletNotFound) {
// 			statusCode = http.StatusPreconditionFailed // Wallet must exist
// 		}
// 		c.JSON(statusCode, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, response)
// }

func (ctrl *CustomerSubscriptionController) GetCustomerTraderSubscriptions(c *gin.Context) {
	customerID := c.MustGet("userID").(uint)

	subscriptions, err := ctrl.TraderSubscriptionService.GetCustomerTraderSubscriptions(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to retrieve subscriptions: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (ctrl *CustomerSubscriptionController) GetCustomerTraderSubscriptionByID(c *gin.Context) {
	customerID := c.MustGet("userID").(uint)
	subIDStr := c.Param("subscription_id")
	subID, err := strconv.ParseUint(subIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid subscription ID format"})
		return
	}

	subscription, err := ctrl.TraderSubscriptionService.GetCustomerTraderSubscriptionByID(c.Request.Context(), customerID, uint(subID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrSubscriptionNotFound) || errors.Is(err, service.ErrNotAuthorized) {
			statusCode = http.StatusNotFound // or Forbidden if it's not their subscription
		}
		c.JSON(statusCode, gin.H{"message": fmt.Sprintf("Failed to retrieve subscription: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, subscription)
}
