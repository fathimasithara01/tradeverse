package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CustomerSubscriptionController struct {
	SubscriptionService service.ITraderSubscriptionService
}

func NewCustomerSubscriptionController(subService service.ITraderSubscriptionService) *CustomerSubscriptionController {
	return &CustomerSubscriptionController{
		SubscriptionService: subService,
	}
}

// GetSubscriptionPlanDetails godoc
// @Summary Get details of a specific trader subscription plan
// @Description Retrieves details of an active trader subscription plan by ID
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription Plan ID"
// @Success 200 {object} models.SubscriptionPlan
// @Failure 400 {object} gin.H "Invalid plan ID"
// @Failure 404 {object} gin.H "Trader subscription plan not found"
// @Failure 500 {object} gin.H "Failed to fetch plan"
// @Router /api/v1/subscriptions/plans/{id} [get]
func (ctrl *CustomerSubscriptionController) GetSubscriptionPlanDetails(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := ctrl.SubscriptionService.GetTraderSubscriptionPlanByID(uint(planID))
	if err != nil {
		if err.Error() == "trader subscription plan not found or not active" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error fetching subscription plan %d: %v", planID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plan"})
		return
	}
	c.JSON(http.StatusOK, plan)
}

// GetAllTraderSubscriptionPlans godoc
// @Summary Get all available trader subscription plans
// @Description Retrieves a list of all active trader subscription plans for customers to choose from
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} models.SubscriptionPlan
// @Failure 500 {object} gin.H "Failed to fetch plans"
// @Router /api/v1/subscriptions/plans [get]
func (ctrl *CustomerSubscriptionController) GetAllTraderSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.SubscriptionService.GetAllActiveTraderSubscriptionPlans()
	if err != nil {
		log.Printf("Error fetching all active trader subscription plans: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// @Router /api/v1/subscriptions [post]
func (ctrl *CustomerSubscriptionController) SubscribeToTrader(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c) // Assuming a helper to get user ID from JWT token
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	var req struct {
		PlanID          uint    `json:"plan_id" binding:"required"`
		TraderProfileID uint    `json:"trader_profile_id" binding:"required"` // The specific trader this subscription is for
		AmountPaid      float64 `json:"amount_paid" binding:"required"`
		TransactionID   string  `json:"transaction_id" binding:"required"` // Payment gateway transaction ID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := ctrl.SubscriptionService.SubscribeToTrader(userID, req.PlanID, req.TraderProfileID, req.AmountPaid, req.TransactionID)
	if err != nil {
		log.Printf("Error creating trader subscription for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetMySubscriptionDetails godoc
// @Summary Get details of a specific customer's trader subscription
// @Description Retrieves details of a specific trader subscription for the authenticated customer.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.TraderSubscription
// @Failure 400 {object} gin.H "Invalid subscription ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Subscription not found or not authorized"
// @Failure 500 {object} gin.H "Failed to fetch subscription"
// @Router /api/v1/subscriptions/{id} [get]
func (ctrl *CustomerSubscriptionController) GetMySubscriptionDetails(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	subIDStr := c.Param("id")
	subID, err := strconv.ParseUint(subIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := ctrl.SubscriptionService.GetMyTraderSubscription(uint(subID), userID)
	if err != nil {
		if err.Error() == "subscription not found or not authorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error fetching subscription %d for user %d: %v", subID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// ListMySubscriptions godoc
// @Summary List all my trader subscriptions
// @Description Retrieves a list of all trader subscriptions for the authenticated customer.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} models.TraderSubscription
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Failed to list subscriptions"
// @Router /api/v1/subscriptions [get]
func (ctrl *CustomerSubscriptionController) ListMySubscriptions(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	subscriptions, err := ctrl.SubscriptionService.ListMyTraderSubscriptions(userID)
	if err != nil {
		log.Printf("Error listing subscriptions for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// UpdateSubscriptionSettings godoc
// @Summary Update trader subscription allocation and risk
// @Description Allows a customer to update the allocation and risk multiplier for an active trader subscription.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param update_request body UpdateSubscriptionSettingsRequest true "New allocation and risk multiplier"
// @Success 200 {object} models.TraderSubscription
// @Failure 400 {object} gin.H "Bad request (validation error, invalid ID)"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Subscription not found or not authorized"
// @Failure 500 {object} gin.H "Failed to update subscription"
// @Router /api/v1/subscriptions/{id} [put]
func (ctrl *CustomerSubscriptionController) UpdateSubscriptionSettings(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	subIDStr := c.Param("id")
	subID, err := strconv.ParseUint(subIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req struct {
		Allocation     float64 `json:"allocation" binding:"required,min=0.01"`
		RiskMultiplier float64 `json:"risk_multiplier" binding:"required,min=0.01"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSub, err := ctrl.SubscriptionService.UpdateTraderSubscriptionSettings(uint(subID), userID, req.Allocation, req.RiskMultiplier)
	if err != nil {
		if err.Error() == "subscription not found or not authorized" || err.Error() == "cannot update an inactive subscription" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error updating subscription %d for user %d: %v", subID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSub)
}

// PauseCopyTrading godoc
// @Summary Pause copy trading for a subscription
// @Description Pauses the copy trading activity for a specific trader subscription.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} gin.H "Message: Copy trading paused successfully"
// @Failure 400 {object} gin.H "Invalid subscription ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Subscription not found or not authorized"
// @Failure 409 {object} gin.H "Subscription already paused"
// @Failure 500 {object} gin.H "Failed to pause copy trading"
// @Router /api/v1/subscriptions/{id}/pause [post]
func (ctrl *CustomerSubscriptionController) PauseCopyTrading(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	subIDStr := c.Param("id")
	subID, err := strconv.ParseUint(subIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	err = ctrl.SubscriptionService.PauseTraderCopyTrading(uint(subID), userID)
	if err != nil {
		if err.Error() == "subscription not found or not authorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "subscription is already paused" || err.Error() == "cannot pause an inactive subscription" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error pausing subscription %d for user %d: %v", subID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pause copy trading: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Copy trading paused successfully"})
}

// ResumeCopyTrading godoc
// @Summary Resume copy trading for a subscription
// @Description Resumes the copy trading activity for a previously paused trader subscription.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} gin.H "Message: Copy trading resumed successfully"
// @Failure 400 {object} gin.H "Invalid subscription ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Subscription not found or not authorized"
// @Failure 409 {object} gin.H "Subscription not paused"
// @Failure 500 {object} gin.H "Failed to resume copy trading"
// @Router /api/v1/subscriptions/{id}/resume [post]
func (ctrl *CustomerSubscriptionController) ResumeCopyTrading(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	subIDStr := c.Param("id")
	subID, err := strconv.ParseUint(subIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	err = ctrl.SubscriptionService.ResumeTraderCopyTrading(uint(subID), userID)
	if err != nil {
		if err.Error() == "subscription not found or not authorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "subscription is not paused" || err.Error() == "cannot resume an inactive subscription" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error resuming subscription %d for user %d: %v", subID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resume copy trading: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Copy trading resumed successfully"})
}

// CancelSubscription godoc
// @Summary Cancel a trader subscription
// @Description Allows a customer to cancel their trader subscription.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} gin.H "Message: Subscription cancelled successfully"
// @Failure 400 {object} gin.H "Invalid subscription ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Subscription not found or not authorized"
// @Failure 409 {object} gin.H "Subscription already inactive"
// @Failure 500 {object} gin.H "Failed to cancel subscription"
// @Router /api/v1/subscriptions/{id} [delete]
func (ctrl *CustomerSubscriptionController) CancelSubscription(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	subIDStr := c.Param("id")
	subID, err := strconv.ParseUint(subIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	err = ctrl.SubscriptionService.CancelTraderSubscription(uint(subID), userID)
	if err != nil {
		if err.Error() == "subscription not found or not authorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "subscription is already inactive or cancelled" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error cancelling subscription %d for user %d: %v", subID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

// SimulateSubscription godoc
// @Summary Run simulation before subscribing
// @Description Runs a simulation for a given trader subscription plan with specified initial capital.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param plan_id path int true "Subscription Plan ID"
// @Param simulation_request body SimulateSubscriptionRequest true "Simulation parameters"
// @Success 200 {object} map[string]interface{} "Simulation result"
// @Failure 400 {object} gin.H "Bad request (validation error, invalid plan ID)"
// @Failure 404 {object} gin.H "Trader subscription plan not found"
// @Failure 500 {object} gin.H "Failed to run simulation"
// @Router /api/v1/subscriptions/plans/{plan_id}/simulate [post]
func (ctrl *CustomerSubscriptionController) SimulateSubscription(c *gin.Context) {
	planIDStr := c.Param("plan_id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var req struct {
		InitialCapital float64 `json:"initial_capital" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.SubscriptionService.SimulateTraderSubscription(uint(planID), req.InitialCapital)
	if err != nil {
		if err.Error() == "trader subscription plan not found for simulation" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error running simulation for plan %d: %v", planID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run simulation: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
