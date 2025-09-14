package controllers

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/fathimasithara01/tradeverse/internal/customer/service"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type TraderSubscriptionController struct {
// 	traderSubscriptionService service.TraderSubscriptionService
// }

// func NewTraderSubscriptionController(traderSubscriptionService service.TraderSubscriptionService) *TraderSubscriptionController {
// 	return &TraderSubscriptionController{traderSubscriptionService: traderSubscriptionService}
// }

// // ListTraderSubscriptionPlans godoc
// // @Summary List all available trader subscription plans
// // @Description Retrieves a list of all active subscription plans for users who want to become traders.
// // @Tags Trader Subscriptions
// // @Produce json
// // @Success 200 {array} models.TraderSubscriptionPlan
// // @Failure 500 {object} map[string]string "Internal Server Error"
// // @Router /api/v1/trader-subscription-plans [get]
// func (c *TraderSubscriptionController) ListTraderSubscriptionPlans(ctx *gin.Context) {
// 	plans, err := c.traderSubscriptionService.ListTraderSubscriptionPlans()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, plans)
// }

// // UpgradeToTrader godoc
// // @Summary Upgrade to a trader account
// // @Description Allows a customer to subscribe to a trader plan, which changes their role to 'trader' and creates a trader profile.
// // @Tags Trader Subscriptions
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param upgradeRequest body UpgradeToTraderRequest true "Upgrade details"
// // @Success 201 {object} models.TraderSubscription
// // @Failure 400 {object} map[string]string "Bad Request"
// // @Failure 401 {object} map[string]string "Unauthorized"
// // @Failure 403 {object} map[string]string "Forbidden - already a trader or invalid role"
// // @Failure 500 {object} map[string]string "Internal Server Error"
// // @Router /api/v1/upgrade-to-trader [post]
// func (c *TraderSubscriptionController) UpgradeToTrader(ctx *gin.Context) {
// 	userID, exists := ctx.Get("userID")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
// 		return
// 	}

// 	var req struct {
// 		TraderSubscriptionPlanID uint   `json:"trader_subscription_plan_id" binding:"required"`
// 		PaymentToken             string `json:"payment_token" binding:"required"` // Token from frontend (e.g., Stripe, PayPal)
// 	}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	traderSubscription, err := c.traderSubscriptionService.UpgradeToTrader(userID.(uint), req.TraderSubscriptionPlanID, req.PaymentToken)
// 	if err != nil {
// 		if err.Error() == "only customers can upgrade to trader" || err.Error() == "user already has an active trader subscription" {
// 			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if err.Error() == "trader subscription plan not found" || err.Error() == "trader subscription plan is not active" {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, traderSubscription)
// }

// // GetMyTraderSubscription godoc
// // @Summary Get current active trader subscription
// // @Description Retrieves the details of the authenticated user's active trader subscription.
// // @Tags Trader Subscriptions
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {object} models.TraderSubscription
// // @Failure 401 {object} map[string]string "Unauthorized"
// // @Failure 404 {object} map[string]string "No active trader subscription found"
// // @Failure 500 {object} map[string]string "Internal Server Error"
// // @Router /api/v1/my-trader-subscription [get]
// func (c *TraderSubscriptionController) GetMyTraderSubscription(ctx *gin.Context) {
// 	userID, exists := ctx.Get("userID")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
// 		return
// 	}

// 	traderSubscription, err := c.traderSubscriptionService.GetMyTraderSubscription(userID.(uint))
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "No active trader subscription found"})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, traderSubscription)
// }

// // This struct is for Swagger documentation.
// type UpgradeToTraderRequest struct {
// 	TraderSubscriptionPlanID uint   `json:"trader_subscription_plan_id" example:"1"`
// 	PaymentToken             string `json:"payment_token" example:"tok_visa"` // A token obtained from a payment gateway
// }
