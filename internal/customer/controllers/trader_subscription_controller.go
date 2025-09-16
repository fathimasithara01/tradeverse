// internal/customer/controllers/customer_controller.go
package controllers

import (
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

func (ctrl *customerController) ListTraderSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.service.ListTraderSubscriptionPlans() // Calls the correct service method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (ctrl *customerController) SubscribeToTraderPlan(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	planIDStr := c.Param("plan_id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid plan ID"})
		return
	}

	sub, err := ctrl.service.SubscribeToTraderPlan(userID, uint(planID))
	if err != nil {
		switch err {
		case service.ErrPlanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		case service.ErrAlreadyHasTraderSubscription, service.ErrNotTraderPlan:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to subscribe: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, sub)
}

func (ctrl *customerController) GetCustomerTraderSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	sub, err := ctrl.service.GetCustomerTraderSubscription(userID)
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

func (ctrl *customerController) CancelCustomerTraderSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	subIDStr := c.Param("subscription_id")
	subID, err := strconv.ParseUint(subIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid subscription ID"})
		return
	}

	err = ctrl.service.CancelCustomerTraderSubscription(userID, uint(subID))
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
