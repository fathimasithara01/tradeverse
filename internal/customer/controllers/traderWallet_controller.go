package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type CustomerWalletController struct {
	svc service.CustomerWalletService
}

func NewCustomerWalletController(svc service.CustomerWalletService) *CustomerWalletController {
	return &CustomerWalletController{svc: svc}
}

func (c *CustomerWalletController) GetBalance(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	res, err := c.svc.GetBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *CustomerWalletController) Deposit(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var input models.DepositRequestInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.svc.Deposit(userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *CustomerWalletController) VerifyDeposit(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var input models.DepositVerifyInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.VerifyDeposit(input.PaymentGatewayTxID, userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Deposit verified & wallet credited"})
}

func (c *CustomerWalletController) Withdraw(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var input models.WithdrawalRequestInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.svc.Withdraw(userID, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
