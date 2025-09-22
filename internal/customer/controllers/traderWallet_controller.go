package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
)

type TraderWalletController struct {
	svc *service.TraderWalletService
}

func NewTraderWalletController(svc *service.TraderWalletService) *TraderWalletController {
	return &TraderWalletController{svc}
}

func (c *TraderWalletController) SubscribeCustomer(ctx *gin.Context) {
	var req struct {
		CustomerID uint    `json:"customer_id" binding:"required"`
		TraderID   uint    `json:"trader_id" binding:"required"`
		Price      float64 `json:"price" binding:"required,gt=0"`
		Currency   string  `json:"currency" binding:"required,oneof=USD INR"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := c.svc.SubscribeCustomer(req.CustomerID, req.TraderID, req.Price, req.Currency); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription successful"})
}

// GetBalance returns wallet balance for the authenticated user
func (c *TraderWalletController) GetBalance(ctx *gin.Context) {
	// Retrieve user_id from JWT claims set by AuthMiddleware
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	balanceResp, err := c.svc.GetBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, balanceResp)
}
