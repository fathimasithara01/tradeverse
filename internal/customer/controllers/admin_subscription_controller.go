// internal/customer/controllers/subscription_plan_controller.go
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionPlanController struct {
	SubscriptionPlanService service.ICustomerSubscriptionPlanService
	SubscriptionService     service.ICustomerSubscriptionService
	WalletService           service.IWalletService
}

func NewSubscriptionPlanController(
	planService service.ICustomerSubscriptionPlanService,
	subService service.ICustomerSubscriptionService,
	walletService service.IWalletService,
) *SubscriptionPlanController {
	return &SubscriptionPlanController{
		SubscriptionPlanService: planService,
		SubscriptionService:     subService,
		WalletService:           walletService,
	}
}

func (ctrl *SubscriptionPlanController) GetAllSubscriptionPlans(c *gin.Context) {
	plans, err := ctrl.SubscriptionPlanService.GetAllSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (ctrl *SubscriptionPlanController) GetSubscriptionPlanByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, plan)
}
func (ctrl *SubscriptionPlanController) SubscribeToPlan(c *gin.Context) {
	userIDValue, exists := c.Get("userID") // match middleware key exactly
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var userID uint
	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	planIDParam := c.Param("id")
	planID, err := strconv.ParseUint(planIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(planID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
		return
	}

	// ✅ Generate a safe transaction ID
	transactionID := fmt.Sprintf("SUB_TX_%d_%d_%d", userID, planID, time.Now().UnixNano())

	err = ctrl.WalletService.DebitUserWallet(userID, plan.Price, plan.Currency, "Subscription to "+plan.Name, transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to debit user wallet: " + err.Error()})
		return
	}

	subscription, err := ctrl.SubscriptionService.CreateSubscription(userID, uint(planID), plan.Price, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Successfully subscribed to plan",
		"subscription":  subscription,
		"transactionID": transactionID,
	})
}

// func (ctrl *SubscriptionPlanController) SubscribeToPlan(c *gin.Context) {
// 	// 1️⃣ Get the user ID from context safely
// 	userIDValue, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
// 		return
// 	}

// 	var userID int64
// 	switch v := userIDValue.(type) {
// 	case int64:
// 		userID = v
// 	case float64:
// 		userID = int64(v)
// 	case int:
// 		userID = int64(v)
// 	case uint:
// 		userID = int64(v)
// 	default:
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
// 		return
// 	}

// 	// 2️⃣ Parse plan ID from URL
// 	planIDParam := c.Param("id")
// 	planID, err := strconv.ParseUint(planIDParam, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
// 		return
// 	}

// 	// 3️⃣ Fetch the subscription plan
// 	plan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(planID))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
// 		return
// 	}

// 	var requestIDStr string
// 	if reqIDVal := c.Request.Context().Value("requestID"); reqIDVal != nil {
// 		switch v := reqIDVal.(type) {
// 		case int64:
// 			requestIDStr = strconv.FormatInt(v, 10)
// 		case string:
// 			requestIDStr = v
// 		default:
// 			requestIDStr = "RANDOM"
// 		}
// 	} else {
// 		requestIDStr = strconv.FormatInt(time.Now().UnixNano(), 10)
// 	}

// 	transactionID := fmt.Sprintf("SUB_TX_%d_%d_%s", userID, planID, requestIDStr)

// 	err = ctrl.WalletService.DebitUserWallet(uint(userID), plan.Price, plan.Currency, "Subscription to "+plan.Name, transactionID)
// 	if err != nil {
// 		log.Printf("Error debiting user wallet: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to debit user wallet: " + err.Error()})
// 		return
// 	}

// 	subscription, err := ctrl.SubscriptionService.CreateSubscription(uint(userID), uint(planID), plan.Price, transactionID)
// 	if err != nil {
// 		log.Printf("Error creating subscription: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription: " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message":       "Successfully subscribed to plan",
// 		"subscription":  subscription,
// 		"transactionID": transactionID,
// 	})
// }

// func (ctrl *SubscriptionPlanController) SubscribeToPlan(c *gin.Context) {

// 	userIDValue, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
// 		return
// 	}
// 	userID, ok := userIDValue.(int64)
// 	if !ok {
// 		// try other possible types (for example, JWT middleware may store as float64)
// 		if f, ok := userIDValue.(float64); ok {
// 			userID = int64(f)
// 		} else if i, ok := userIDValue.(int); ok {
// 			userID = int64(i)
// 		} else if u, ok := userIDValue.(uint); ok {
// 			userID = int64(u)
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
// 			return
// 		}
// 	}

// 	planIDParam := c.Param("id")
// 	planID, err := strconv.ParseUint(planIDParam, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
// 		return
// 	}

// 	plan, err := ctrl.SubscriptionPlanService.GetSubscriptionPlanByID(uint(planID))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
// 		return
// 	}

// 	transactionID := "SUB_TX_" + strconv.FormatUint(uint64(userID), 10) + "_" + strconv.FormatUint(uint64(planID), 10) + "_" + strconv.FormatInt(c.Request.Context().Value("requestID").(int64), 10) // Example unique ID

// 	err = ctrl.WalletService.DebitUserWallet(userID, plan.Price, plan.Currency, "Subscription to "+plan.Name, transactionID)
// 	if err != nil {
// 		log.Printf("Error debiting user wallet: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to debit user wallet: " + err.Error()})
// 		return
// 	}

// 	subscription, err := ctrl.SubscriptionService.CreateSubscription(userID, uint(planID), plan.Price, transactionID)
// 	if err != nil {
// 		log.Printf("Error creating subscription: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription: " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Successfully subscribed to plan", "subscription": subscription})
// }

func (ctrl *SubscriptionPlanController) CancelSubscription(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	subscriptionIDParam := c.Param("id")
	subscriptionID, err := strconv.ParseUint(subscriptionIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	err = ctrl.SubscriptionService.CancelSubscription(userID, uint(subscriptionID))
	if err != nil {
		log.Printf("Error cancelling subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

func (ctrl *SubscriptionPlanController) GetUserSubscriptions(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	subscriptions, err := ctrl.SubscriptionService.GetSubscriptionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user subscriptions"})
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}
