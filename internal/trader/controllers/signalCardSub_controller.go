// internal/trader/controllers/trader_subscription.go
package controllers

import (
	"fmt"
	"log"
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

func getUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID.(uint), nil
}

func (ctrl *TraderSubscriptionController) CreateTraderSubscriptionPlan(c *gin.Context) {
	traderID, err := getUserID(c)
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

func (ctrl *TraderSubscriptionController) GetMyTraderSubscriptionPlans(c *gin.Context) {
	traderID, err := getUserID(c)
	if err != nil {
		log.Printf("ERROR: GetMyTraderSubscriptionPlans - Unauthorized: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	log.Printf("INFO: GetMyTraderSubscriptionPlans - Fetching plans for trader ID: %d", traderID)

	plans, err := ctrl.subsService.GetMyTraderSubscriptionPlans(c, traderID)
	if err != nil {
		log.Printf("ERROR: GetMyTraderSubscriptionPlans - Service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch trader's plans: %v", err)})
		return
	}
	log.Printf("INFO: GetMyTraderSubscriptionPlans - Found %d plans for trader ID %d", len(plans), traderID)
	c.JSON(http.StatusOK, plans)
}

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

func (ctrl *TraderSubscriptionController) UpdateTraderSubscriptionPlan(c *gin.Context) {
	traderID, err := getUserID(c)
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

func (ctrl *TraderSubscriptionController) DeleteTraderSubscriptionPlan(c *gin.Context) {
	traderID, err := getUserID(c)
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
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()}) 
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete plan: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trader subscription plan deleted successfully"})
}

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

func (ctrl *TraderSubscriptionController) GetAllTraderUpgradePlans(c *gin.Context) {
	plans, err := ctrl.subsService.GetAllTraderUpgradePlans(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch trader upgrade plans: %v", err)})
		return
	}
	c.JSON(http.StatusOK, plans)
}

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
