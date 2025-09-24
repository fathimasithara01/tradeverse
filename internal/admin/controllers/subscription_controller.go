package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Re-using the DTO from before, but will apply it for input as well
type SubscriptionPlanResponseDTO struct {
	ID              uint    `json:"ID"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Duration        int     `json:"duration"`
	Interval        string  `json:"interval"`
	MaxFollowers    int     `json:"max_followers"`
	Status          string  `json:"status"` // Calculated status for display
	Features        string  `json:"features"`
	CommissionRate  float64 `json:"commission_rate"`
	AnalyticsAccess string  `json:"analytics_access"`
	IsTraderPlan    bool    `json:"is_trader_plan"` // Raw boolean for form population
	IsActive        bool    `json:"is_active"`      // Raw boolean for form population
}

type CreateUpdateSubscriptionPlanRequest struct {
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" binding:"required,gt=0"`
	Duration        int     `json:"duration" binding:"required,gt=0"`
	Interval        string  `json:"interval"`
	MaxFollowers    int     `json:"max_followers"`
	Features        string  `json:"features"`
	CommissionRate  float64 `json:"commission_rate"`
	AnalyticsAccess string  `json:"analytics_access"`
	IsTraderPlan    bool    `json:"is_trader_plan"` // Admin can specify if it's a trader plan
	IsActive        bool    `json:"is_active"`      // Admin can specify if it's active
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

// ShowEditSubscriptionPlanPage is likely not needed with the current frontend approach
// The frontend fetches data via API and populates a modal.
/*
func (ctrl *SubscriptionController) ShowEditSubscriptionPlanPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid plan ID"})
		return
	}
	plan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Subscription plan not found"})
			return
		}
		log.Printf("Error fetching plan for edit page: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to retrieve subscription plan"})
		return
	}
	c.HTML(http.StatusOK, "edit_subscription_plan.html", gin.H{
		"Title":  "Edit Subscription Plan",
		"ActiveTab":    "financials",
		"ActiveSubTab": "subscription_plans",
		"Plan": plan, // Pass the entire plan object if the page truly expects it
	})
}
*/

// GetSubscriptionPlanByID API endpoint for the frontend to fetch plan data for editing
func (ctrl *SubscriptionController) GetSubscriptionPlanByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
			return
		}
		log.Printf("Error fetching subscription plan by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription plan"})
		return
	}

	// Prepare response DTO
	status := "inactive"
	if plan.IsActive {
		status = "active"
	}
	responsePlan := SubscriptionPlanResponseDTO{
		ID:              plan.ID,
		Name:            plan.Name,
		Description:     plan.Description,
		Price:           plan.Price,
		Duration:        plan.Duration,
		Interval:        plan.Interval,
		MaxFollowers:    plan.MaxFollowers,
		Status:          status,
		Features:        plan.Features,
		CommissionRate:  plan.CommissionRate,
		AnalyticsAccess: plan.AnalyticsAccess,
		IsTraderPlan:    plan.IsTraderPlan,
		IsActive:        plan.IsActive,
	}

	c.JSON(http.StatusOK, responsePlan)
}

func (ctrl *SubscriptionController) GetSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.SubscriptionPlanService.GetAllSubscriptionPlans()
	if err != nil {
		log.Printf("Error fetching subscription plans: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription plans"})
		return
	}

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
			IsActive:        plan.IsActive, // Expose IsActive in response
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
	var req CreateUpdateSubscriptionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	newPlan := models.SubscriptionPlan{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Duration:        req.Duration,
		Interval:        req.Interval,
		MaxFollowers:    req.MaxFollowers,
		Features:        req.Features,
		CommissionRate:  req.CommissionRate,
		AnalyticsAccess: req.AnalyticsAccess,
		IsTraderPlan:    req.IsTraderPlan,
		IsActive:        req.IsActive,
	}

	if err := ctrl.SubscriptionPlanService.CreateSubscriptionPlan(&newPlan); err != nil {
		log.Printf("Error creating subscription plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to create subscription plan: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SubscriptionPlanResponseDTO{
		ID:              newPlan.ID,
		Name:            newPlan.Name,
		Description:     newPlan.Description,
		Price:           newPlan.Price,
		Duration:        newPlan.Duration,
		Interval:        newPlan.Interval,
		MaxFollowers:    newPlan.MaxFollowers,
		Status:          map[bool]string{true: "active", false: "inactive"}[newPlan.IsActive],
		Features:        newPlan.Features,
		CommissionRate:  newPlan.CommissionRate,
		AnalyticsAccess: newPlan.AnalyticsAccess,
		IsTraderPlan:    newPlan.IsTraderPlan,
		IsActive:        newPlan.IsActive,
	})
}

func (ctrl *SubscriptionController) UpdateSubscriptionPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var req CreateUpdateSubscriptionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingPlan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
			return
		}
		log.Printf("Error fetching existing plan for update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription plan"})
		return
	}

	existingPlan.Name = req.Name
	existingPlan.Description = req.Description
	existingPlan.Price = req.Price
	existingPlan.Duration = req.Duration
	existingPlan.Interval = req.Interval
	existingPlan.MaxFollowers = req.MaxFollowers
	existingPlan.Features = req.Features
	existingPlan.CommissionRate = req.CommissionRate
	existingPlan.AnalyticsAccess = req.AnalyticsAccess
	existingPlan.IsTraderPlan = req.IsTraderPlan
	existingPlan.IsActive = req.IsActive

	if err := ctrl.SubscriptionPlanService.UpdateSubscriptionPlan(existingPlan); err != nil {
		log.Printf("Error updating subscription plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription plan"})
		return
	}

	status := "inactive"
	if existingPlan.IsActive {
		status = "active"
	}

	responsePlan := SubscriptionPlanResponseDTO{
		ID:              existingPlan.ID,
		Name:            existingPlan.Name,
		Description:     existingPlan.Description,
		Price:           existingPlan.Price,
		Duration:        existingPlan.Duration,
		Interval:        existingPlan.Interval,
		MaxFollowers:    existingPlan.MaxFollowers,
		Status:          status,
		Features:        existingPlan.Features,
		CommissionRate:  existingPlan.CommissionRate,
		AnalyticsAccess: existingPlan.AnalyticsAccess,
		IsTraderPlan:    existingPlan.IsTraderPlan,
		IsActive:        existingPlan.IsActive,
	}

	c.JSON(http.StatusOK, responsePlan)
}

func (ctrl *SubscriptionController) DeleteSubscriptionPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	if err := ctrl.SubscriptionPlanService.DeleteSubscriptionPlan(uint(id)); err != nil {
		log.Printf("Error deleting subscription plan %d: %v", id, err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription plan deleted successfully"})
}
