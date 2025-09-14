package controllers

import (
	"log"
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

func (ctrl *SubscriptionController) ShowSubscriptionsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_subscriptions.html", gin.H{
		"Title":        "Customer Subscriptions",
		"ActiveTab":    "financials",
		"ActiveSubTab": "subscriptions",
	})
}

func (ctrl *SubscriptionController) ShowSubscriptionPlansPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_subscription_plans.html", gin.H{
		"Title":        "Subscription Plans Management",
		"ActiveTab":    "financials",
		"ActiveSubTab": "subscription_plans",
	})
}

func (ctrl *SubscriptionController) GetSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.SubscriptionPlanService.GetAllSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (ctrl *SubscriptionController) GetSubscriptions(c *gin.Context) {
	subscriptions, err := ctrl.SubscriptionService.GetAllSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}

func (ctrl *SubscriptionController) CreateCustomerSubscription(c *gin.Context) {
	var req struct {
		UserID          uint    `json:"user_id" binding:"required"`
		PlanID          uint    `json:"plan_id" binding:"required"`
		AmountPaid      float64 `json:"amount_paid" binding:"required"`
		TransactionID   string  `json:"transaction_id" binding:"required"`
		IsTraderUpgrade bool    `json:"is_trader_upgrade"`
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

	if req.IsTraderUpgrade {
		err := ctrl.SubscriptionService.UpgradeUserToTrader(req.UserID)
		if err != nil {
			log.Printf("Warning: Failed to upgrade user %d to trader role after subscription: %v", req.UserID, err)
		}
	}

	log.Printf("Payment received for plan %d from user %d: $%.2f. Transaction ID: %s. (To be deposited into admin wallet)",
		req.PlanID, req.UserID, req.AmountPaid, req.TransactionID)

	c.JSON(http.StatusCreated, subscription)
}

func (ctrl *SubscriptionController) CreateSubscriptionPlan(c *gin.Context) {
	var newPlan models.SubscriptionPlan
	if err := c.ShouldBindJSON(&newPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.SubscriptionPlanService.CreateSubscriptionPlan(&newPlan); err != nil {
		log.Printf("Error creating subscription plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription plan"})
		return
	}
	c.JSON(http.StatusCreated, newPlan)
}

func (ctrl *SubscriptionController) UpdateSubscriptionPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var updatedPlanData models.SubscriptionPlan
	if err := c.ShouldBindJSON(&updatedPlanData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlanData.ID = uint(id) // Ensure the ID from the URL is used for the update

	if err := ctrl.SubscriptionPlanService.UpdateSubscriptionPlan(&updatedPlanData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription plan"})
		return
	}
	c.JSON(http.StatusOK, updatedPlanData) // Return the updated plan
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

// func (ctrl *SubscriptionController) CreateCustomerSubscription(c *gin.Context) {
// 	var req struct {
// 		UserID        uint    `json:"user_id" binding:"required"`
// 		PlanID        uint    `json:"plan_id" binding:"required"`
// 		AmountPaid    float64 `json:"amount_paid" binding:"required"`
// 		TransactionID string  `json:"transaction_id" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	subscription, err := ctrl.SubscriptionService.CreateSubscription(req.UserID, req.PlanID, req.AmountPaid, req.TransactionID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer subscription: " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, subscription)
// }
