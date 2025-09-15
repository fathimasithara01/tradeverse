package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
)

type CustomerController interface {
	ListTraderSubscriptionPlans(c *gin.Context)
	SubscribeToTraderPlan(c *gin.Context)
	GetCustomerTraderSubscription(c *gin.Context)
	CancelCustomerTraderSubscription(c *gin.Context)
}

type customerController struct {
	service service.CustomerService
}

func NewCustomerController(s service.CustomerService) CustomerController {
	return &customerController{service: s}
}

// @Router /customer/trader-plans [get]
func (ctrl *customerController) ListTraderSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.service.ListTraderSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// @Router /customer/trader-plans/{plan_id}/subscribe [post]
func (ctrl *customerController) SubscribeToTraderPlan(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Get userID from JWT claims
	planIDStr := c.Param("plan_id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid plan ID"})
		return
	}

	sub, err := ctrl.service.SubscribeToTraderPlan(userID, uint(planID))
	if err != nil {
		if errors.Is(err, errors.New("subscription plan not found")) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, errors.New("user already has an active trader subscription")) || errors.Is(err, errors.New("this is not a trader subscription plan")) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to subscribe: %v", err)})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// @Summary Get customer's active Trader Subscription
// @Description Get the details of the currently active trader subscription for the logged-in customer.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.UserTraderSubscriptionResponse "Returns active subscription"
// @Failure 404 {object} map[string]string "message: No active trader subscription found"
// @Failure 500 {object} map[string]string "message: Failed to fetch subscription"
// @Router /customer/trader-subscription [get]
func (ctrl *customerController) GetCustomerTraderSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Get userID from JWT claims

	sub, err := ctrl.service.GetCustomerTraderSubscription(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to retrieve subscription: %v", err)})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No active trader subscription found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// @Router /customer/trader-subscription/{subscription_id}/cancel [post]
func (ctrl *customerController) CancelCustomerTraderSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Get userID from JWT claims
	subIDStr := c.Param("subscription_id")
	subID, err := strconv.ParseUint(subIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid subscription ID"})
		return
	}

	err = ctrl.service.CancelCustomerTraderSubscription(userID, uint(subID))
	if err != nil {
		if errors.Is(err, errors.New("active trader subscription not found for this user and ID")) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to cancel subscription: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trader subscription cancelled successfully"})
}
