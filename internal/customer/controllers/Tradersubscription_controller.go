package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	svc service.SubscriptionService
}

func NewSubscriptionController(svc service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{svc: svc}
}

func (c *SubscriptionController) SubscribeCustomer(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	planID, err := strconv.Atoi(ctx.Param("plan_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	sub, err := c.svc.SubscribeCustomerToTrader(userID, uint(planID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "subscription successful",
		"subscription": sub,
	})
}
