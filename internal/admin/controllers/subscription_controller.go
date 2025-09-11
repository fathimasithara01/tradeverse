package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	SubscriptionService     service.ISubscriptionService
	SubscriptionPlanService service.ISubscriptionPlanService
}

func NewSubscriptionController(subService service.ISubscriptionService, planService service.ISubscriptionPlanService) *SubscriptionController {
	return &SubscriptionController{
		SubscriptionService:     subService,
		SubscriptionPlanService: planService,
	}
}

// ShowSubscriptionsPage renders the HTML page for viewing all customer subscriptions
func (ctrl *SubscriptionController) ShowSubscriptionsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_subscriptions.html", gin.H{
		"Title":        "Customer Subscriptions",
		"ActiveTab":    "financials",
		"ActiveSubTab": "subscriptions",
	})
}

// GetSubscriptions fetches all customer subscriptions with associated user and plan details
func (ctrl *SubscriptionController) GetSubscriptions(c *gin.Context) {
	subscriptions, err := ctrl.SubscriptionService.GetAllSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}

// ShowSubscriptionPlansPage renders the HTML page for managing subscription plans
func (ctrl *SubscriptionController) ShowSubscriptionPlansPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_subscription_plans.html", gin.H{
		"Title":        "Subscription Plans Management",
		"ActiveTab":    "financials",
		"ActiveSubTab": "subscription_plans",
	})
}

// GetSubscriptionPlans fetches all subscription plans
func (ctrl *SubscriptionController) GetSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.SubscriptionPlanService.GetAllSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// CreateSubscriptionPlan creates a new subscription plan
func (ctrl *SubscriptionController) CreateSubscriptionPlan(c *gin.Context) {
	var plan models.SubscriptionPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.SubscriptionPlanService.CreateSubscriptionPlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription plan"})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

// UpdateSubscriptionPlan updates an existing subscription plan
func (ctrl *SubscriptionController) UpdateSubscriptionPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
		return
	}

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID is preserved from the URL parameter for the update
	plan.ID = uint(id)

	if err := ctrl.SubscriptionPlanService.UpdateSubscriptionPlan(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription plan"})
		return
	}
	c.JSON(http.StatusOK, plan)
}

// DeleteSubscriptionPlan deletes a subscription plan
func (ctrl *SubscriptionController) DeleteSubscriptionPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	if err := ctrl.SubscriptionPlanService.DeleteSubscriptionPlan(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription plan deleted successfully"})
}

// CreateCustomerSubscription handles creating a subscription for a customer
// This would typically be called by a webhook from a payment gateway, not directly from admin UI.
func (ctrl *SubscriptionController) CreateCustomerSubscription(c *gin.Context) {
	var req struct {
		UserID        uint    `json:"user_id" binding:"required"`
		PlanID        uint    `json:"plan_id" binding:"required"`
		AmountPaid    float64 `json:"amount_paid" binding:"required"`
		TransactionID string  `json:"transaction_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := ctrl.SubscriptionService.CreateSubscription(req.UserID, req.PlanID, req.AmountPaid, req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}
