// internal/trader/controllers/trader_subscription.go
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type TraderSubscriptionController struct {
	subsService service.ITraderSubscriptionService
}

func NewTraderSubscriptionController(subsService service.ITraderSubscriptionService) *TraderSubscriptionController {
	return &TraderSubscriptionController{subsService: subsService}
}

// Helper to get user ID from context (can be customer or trader)
func getUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID.(uint), nil
}

// @Summary Create a new subscription plan for the current trader
// @Description Allows an authenticated trader to create a new subscription plan for their signals.
// @Tags Trader Subscriptions
// @Accept json
// @Produce json
// @Param planInput body models.CreateTraderSubscriptionPlanInput true "Trader Subscription Plan details"
// @Success 201 {object} models.TraderSubscriptionPlan "Successfully created plan"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (user is not an active trader)"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /trader/plans [post]
func (ctrl *TraderSubscriptionController) CreateTraderSubscriptionPlan(c *gin.Context) {
	traderID, err := getUserID(c) // Use general getUserID
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input models.CreateTraderSubscriptionPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := ctrl.subsService.CreateTraderSubscriptionPlan(c, traderID, input)
	if err != nil {
		if err.Error() == "user is not an active trader and cannot create subscription plans" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create subscription plan: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// @Summary Get all subscription plans created by the current trader
// @Description Retrieve a list of all subscription plans created by the authenticated trader.
// @Tags Trader Subscriptions
// @Produce json
// @Success 200 {array} models.TraderSubscriptionPlan "List of trader's plans"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /trader/plans [get]
func (ctrl *TraderSubscriptionController) GetMyTraderSubscriptionPlans(c *gin.Context) {
	traderID, err := getUserID(c) // Use general getUserID
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	plans, err := ctrl.subsService.GetMyTraderSubscriptionPlans(c, traderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch trader's plans: %v", err)})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// @Summary Get a specific trader subscription plan by ID
// @Description Retrieve details of a specific subscription plan by its ID.
// @Tags Trader Subscriptions
// @Produce json
// @Param planId path int true "Plan ID"
// @Success 200 {object} models.TraderSubscriptionPlan "Trader subscription plan details"
// @Failure 400 {object} gin.H "Invalid plan ID"
// @Failure 404 {object} gin.H "Plan not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /trader/plans/:planId [get]
func (ctrl *TraderSubscriptionController) GetTraderSubscriptionPlanByID(c *gin.Context) {
	planIDParam := c.Param("planId")
	planID, err := strconv.ParseUint(planIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	plan, err := ctrl.subsService.GetTraderSubscriptionPlanByID(c, uint(planID))
	if err != nil {
		if err.Error() == "trader subscription plan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get plan: %v", err)})
		return
	}
	c.JSON(http.StatusOK, plan)
}

// @Summary Update a trader subscription plan
// @Description Allows an authenticated trader to update their existing subscription plan.
// @Tags Trader Subscriptions
// @Accept json
// @Produce json
// @Param planId path int true "Plan ID"
// @Param planInput body models.CreateTraderSubscriptionPlanInput true "Updated Trader Subscription Plan details"
// @Success 200 {object} models.TraderSubscriptionPlan "Successfully updated plan"
// @Failure 400 {object} gin.H "Bad request or invalid plan ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (if plan does not belong to trader)"
// @Failure 404 {object} gin.H "Plan not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /trader/plans/:planId [put]
func (ctrl *TraderSubscriptionController) UpdateTraderSubscriptionPlan(c *gin.Context) {
	traderID, err := getUserID(c) // Use general getUserID
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	planIDParam := c.Param("planId")
	planID, err := strconv.ParseUint(planIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var input models.CreateTraderSubscriptionPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := ctrl.subsService.UpdateTraderSubscriptionPlan(c, traderID, uint(planID), input)
	if err != nil {
		if err.Error() == "unauthorized: plan does not belong to this trader" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "trader subscription plan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update plan: %v", err)})
		return
	}
	c.JSON(http.StatusOK, plan)
}

// @Summary Delete a trader subscription plan
// @Description Allows an authenticated trader to delete one of their subscription plans.
// @Tags Trader Subscriptions
// @Produce json
// @Param planId path int true "Plan ID"
// @Success 200 {object} gin.H "Successfully deleted plan"
// @Failure 400 {object} gin.H "Invalid plan ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (if plan does not belong to trader)"
// @Failure 404 {object} gin.H "Plan not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /trader/plans/:planId [delete]
func (ctrl *TraderSubscriptionController) DeleteTraderSubscriptionPlan(c *gin.Context) {
	traderID, err := getUserID(c) // Use general getUserID
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	planIDParam := c.Param("planId")
	planID, err := strconv.ParseUint(planIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	err = ctrl.subsService.DeleteTraderSubscriptionPlan(c, traderID, uint(planID))
	if err != nil {
		if err.Error() == "trader subscription plan not found or not owned by this trader" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()}) // Use forbidden as it implies not owned
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete plan: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trader subscription plan deleted successfully"})
}

// --- Customer-facing routes to subscribe to a Trader's Plan ---

// @Summary Subscribe to a specific trader's subscription plan
// @Description Allows an authenticated customer to subscribe to a specific trader's plan.
// @Tags Customer Subscriptions
// @Accept json
// @Produce json
// @Param traderId path int true "ID of the Trader whose plan is being subscribed to"
// @Param planId path int true "ID of the Trader Subscription Plan"
// @Success 200 {object} gin.H "Subscription successful"
// @Failure 400 {object} gin.H "Bad request (invalid ID, insufficient funds, plan inactive, already subscribed)"
// @Failure 401 {object} gin.H "Unauthorized (user not logged in)"
// @Failure 404 {object} gin.H "Plan or Trader not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /customer/subscribe/trader/:traderId/plan/:planId [post]
func (ctrl *TraderSubscriptionController) SubscribeToTraderPlan(c *gin.Context) {
	customerID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	traderIDParam := c.Param("traderId")
	traderID, err := strconv.ParseUint(traderIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trader ID"})
		return
	}

	planIDParam := c.Param("planId")
	planID, err := strconv.ParseUint(planIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	err = ctrl.subsService.SubscribeToTraderPlan(c, customerID, uint(traderID), uint(planID))
	if err != nil {
		switch err.Error() {
		case "trader subscription plan not found", "trader not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case "customer cannot subscribe to their own plan", "trader subscription plan is not active", "insufficient funds in wallet", "user is already subscribed to this plan":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to subscribe to trader plan: %v", err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully subscribed to trader's plan!"})
}

// --- Admin Subscription Routes (for a user to become a trader) ---
// These routes would typically be in an Admin or a shared user controller,
// but placed here for context as it affects a user becoming a trader.

// @Summary Get all admin-defined subscription plans (customer to trader upgrade)
// @Description Allows users to see plans to become a trader.
// @Tags Trader Upgrade Subscriptions
// @Produce json
// @Success 200 {array} models.SubscriptionPlan "List of admin subscription plans for upgrading to trader"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /upgrade/plans [get]
func (ctrl *TraderSubscriptionController) GetAllTraderUpgradePlans(c *gin.Context) {
	plans, err := ctrl.subsService.GetAllTraderUpgradePlans(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch trader upgrade plans: %v", err)})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// @Summary Subscribe to an admin plan to become a trader
// @Description Allows a user to subscribe to an admin-defined plan to gain trader privileges.
// @Tags Trader Upgrade Subscriptions
// @Accept json
// @Produce json
// @Param planId path int true "Admin Subscription Plan ID"
// @Success 200 {object} gin.H "Subscription successful, user is now a trader"
// @Failure 400 {object} gin.H "Bad request or insufficient funds"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /upgrade/subscribe/:planId [post]
func (ctrl *TraderSubscriptionController) SubscribeToTraderUpgradePlan(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	planIDParam := c.Param("planId")
	planID, err := strconv.ParseUint(planIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	err = ctrl.subsService.SubscribeToTraderUpgradePlan(c, userID, uint(planID))
	if err != nil {
		if err.Error() == "insufficient funds in wallet" ||
			err.Error() == "subscription plan is not active" ||
			err.Error() == "this plan is not for upgrading to a trader role" ||
			err.Error() == "user is already an active trader with this upgrade plan" { // Added check
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to subscribe to trader upgrade plan: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully subscribed to trader upgrade plan, you are now a trader!"})
}
