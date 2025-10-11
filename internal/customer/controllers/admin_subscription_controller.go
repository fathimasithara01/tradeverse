package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
)

type AdminSubscriptionController struct {
	Service service.AdminSubscriptionService
}

func NewAdminSubscriptionController(service service.AdminSubscriptionService) *AdminSubscriptionController {
	return &AdminSubscriptionController{Service: service}
}

func (ctrl *AdminSubscriptionController) ListTraderSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.Service.ListTraderSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (ctrl *AdminSubscriptionController) SubscribeToTraderPlan(c *gin.Context) {
	var req struct {
		CustomerID uint `json:"customer_id" binding:"required"`
		PlanID     uint `json:"trader_subscription_plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request: " + err.Error()})
		return
	}

	sub, err := ctrl.Service.SubscribeToTraderPlan(req.CustomerID, req.PlanID)
	if err != nil {
		switch err {
		case service.ErrPlanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case service.ErrAlreadyHasTraderSubscription, service.ErrNotTraderPlan:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case service.ErrInsufficientFunds:
			c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient funds"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (ctrl *AdminSubscriptionController) GetCustomerTraderSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	sub, err := ctrl.Service.GetCustomerTraderSubscription(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve subscription: " + err.Error()})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no active trader subscription found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}
func (ctrl *AdminSubscriptionController) CancelCustomerTraderSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	subIDStr := c.Param("subscription_id")
	subID, err := strconv.ParseUint(subIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid subscription ID"})
		return
	}

	err = ctrl.Service.CancelCustomerTraderSubscription(c, userID, uint(subID))
	if err != nil {
		if err == service.ErrNoActiveTraderSubscription {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to cancel subscription: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trader subscription cancelled successfully"})
}
