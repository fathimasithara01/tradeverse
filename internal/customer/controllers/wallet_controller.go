package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type WalletController struct {
	WalletSvc service.WalletService
}

func NewWalletController(walletSvc service.WalletService) *WalletController {
	return &WalletController{WalletSvc: walletSvc}
}

func (ctrl *WalletController) GetWalletSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	summary, err := ctrl.WalletSvc.GetWalletSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (ctrl *WalletController) InitiateDeposit(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var input models.DepositRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resp, err := ctrl.WalletSvc.InitiateDeposit(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) VerifyDeposit(c *gin.Context) {
	depositIDStr := c.Param("deposit_id")
	depositID, err := strconv.ParseUint(depositIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid deposit ID"})
		return
	}

	var input models.DepositVerifyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resp, err := ctrl.WalletSvc.VerifyDeposit(uint(depositID), input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrDepositAlreadyProcessed) || errors.Is(err, service.ErrInvalidDepositStatus) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, service.ErrUserWalletNotFound) || errors.Is(err, repository.ErrDepositRequestNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) RequestWithdrawal(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var input models.WithdrawalRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resp, err := ctrl.WalletSvc.RequestWithdrawal(userID, input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, repository.ErrInsufficientFunds) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, service.ErrUserWalletNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) GetWalletTransactions(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resp, err := ctrl.WalletSvc.GetTransactions(userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
