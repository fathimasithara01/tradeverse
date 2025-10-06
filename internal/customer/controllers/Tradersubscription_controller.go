package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
)

type TraderSubscriptionController struct {
	TraderSubscriptionSvc service.ITraderSubscriptionService
}

func NewTraderSubscriptionController(svc service.ITraderSubscriptionService) *TraderSubscriptionController {
	return &TraderSubscriptionController{TraderSubscriptionSvc: svc}
}

func (ctrl *TraderSubscriptionController) SubscribeToTrader(c *gin.Context) {
	customerID := c.MustGet("userID").(uint)

	traderIDStr := c.Param("trader_id")
	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid trader ID"})
		return
	}

	planIDStr := c.Param("plan_id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid plan ID"})
		return
	}

	subscription, err := ctrl.TraderSubscriptionSvc.SubscribeToTrader(c.Request.Context(), customerID, uint(traderID), uint(planID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrAlreadySubscribed) || errors.Is(err, service.ErrSelfSubscription) ||
			errors.Is(err, service.ErrPlanNotForTrader) || errors.Is(err, service.ErrTraderNotMatchPlan) ||
			errors.Is(err, service.ErrInsufficientFunds) {
			statusCode = http.StatusBadRequest
			if errors.Is(err, service.ErrAlreadySubscribed) {
				statusCode = http.StatusConflict
			}
		} else if errors.Is(err, customerrepo.ErrSubscriptionNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (ctrl *TraderSubscriptionController) GetMyTraderSubscriptions(c *gin.Context) {
	customerID := c.MustGet("userID").(uint)

	subscriptions, err := ctrl.TraderSubscriptionSvc.GetCustomerTraderSubscriptions(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (ctrl *TraderSubscriptionController) GetTraderSubscriptionPlans(c *gin.Context) {
	traderIDStr := c.Param("trader_id")
	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid trader ID"})
		return
	}

	plans, err := ctrl.TraderSubscriptionSvc.GetTraderPlans(c.Request.Context(), uint(traderID))
	if err != nil {
		if errors.Is(err, customerrepo.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Trader not found or has no active plans"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}
