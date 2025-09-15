package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

// SubscriptionPlanResponseDTO is a DTO for returning subscription plans to the frontend,
// including a 'status' string based on IsActive.
type SubscriptionPlanResponseDTO struct {
	ID              uint    `json:"ID"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Duration        int     `json:"duration"`
	Interval        string  `json:"interval"`
	MaxFollowers    int     `json:"max_followers"`
	Status          string  `json:"status"` // "active" or "inactive" based on IsActive
	Features        string  `json:"features"`
	CommissionRate  float64 `json:"commission_rate"`
	AnalyticsAccess string  `json:"analytics_access"`
	IsTraderPlan    bool    `json:"is_trader_plan"`
}

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
		log.Printf("Error fetching subscription plans: %v", err) // Log the actual error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription plans"})
		return
	}

	// Convert models.SubscriptionPlan to SubscriptionPlanResponseDTO
	var responsePlans []SubscriptionPlanResponseDTO
	for _, plan := range plans {
		status := "inactive"
		if plan.IsActive {
			status = "active"
		}
		responsePlans = append(responsePlans, SubscriptionPlanResponseDTO{
			ID:              plan.ID,
			Name:            plan.Name,
			Description:     plan.Description,
			Price:           plan.Price,
			Duration:        plan.Duration,
			Interval:        plan.Interval,
			MaxFollowers:    plan.MaxFollowers,
			Status:          status, // Set status based on IsActive
			Features:        plan.Features,
			CommissionRate:  plan.CommissionRate,
			AnalyticsAccess: plan.AnalyticsAccess,
			IsTraderPlan:    plan.IsTraderPlan,
		})
	}
	c.JSON(http.StatusOK, responsePlans)
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

	// Frontend sends 'status' string, but model expects 'IsActive' boolean.
	// Convert it here. Assume "active" means true, anything else means false.
	// This requires the 'status' field in the incoming JSON to be handled.
	// We need to adjust the frontend to send IsActive directly, or parse a status field here.
	// Let's assume the frontend will send 'IsActive' or we need to add a 'Status' field to the `newPlan` struct and convert.
	// For now, I'll update the frontend's JSON sending logic.

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

	var updatedPlanData models.SubscriptionPlan // Changed to models.SubscriptionPlan
	if err := c.ShouldBindJSON(&updatedPlanData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlanData.ID = uint(id) // Ensure the ID from the URL is used for the update

	// Fetch existing plan to preserve fields not sent in update (e.g., CreatedAt) if needed
	// Or ensure frontend sends all fields required for update.
	// For simplicity, we directly save updatedPlanData assuming it contains all necessary fields or GORM handles zero values.

	if err := ctrl.SubscriptionPlanService.UpdateSubscriptionPlan(&updatedPlanData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription plan"})
		return
	}

	// Convert the updated plan back to DTO for consistent response
	status := "inactive"
	if updatedPlanData.IsActive {
		status = "active"
	}
	responsePlan := SubscriptionPlanResponseDTO{
		ID:              updatedPlanData.ID,
		Name:            updatedPlanData.Name,
		Description:     updatedPlanData.Description,
		Price:           updatedPlanData.Price,
		Duration:        updatedPlanData.Duration,
		Interval:        updatedPlanData.Interval,
		MaxFollowers:    updatedPlanData.MaxFollowers,
		Status:          status,
		Features:        updatedPlanData.Features,
		CommissionRate:  updatedPlanData.CommissionRate,
		AnalyticsAccess: updatedPlanData.AnalyticsAccess,
		IsTraderPlan:    updatedPlanData.IsTraderPlan,
	}

	c.JSON(http.StatusOK, responsePlan) // Return the updated plan as DTO
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
