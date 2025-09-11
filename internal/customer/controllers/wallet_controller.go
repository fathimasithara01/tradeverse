package controllers

import (
	"errors" // Import errors for error checking
	"fmt"
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type WalletController struct {
	WalletSvc service.WalletServicer
}

func NewWalletController(walletSvc service.WalletServicer) *WalletController {
	return &WalletController{WalletSvc: walletSvc}
}

func (ctrl *WalletController) GetWalletSummary(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication required."})
		return
	}

	summary, err := ctrl.WalletSvc.GetWalletSummary(userID.(uint))
	if err != nil {
		// Use errors.Is for specific known errors if needed, otherwise generic 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallet summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *WalletController) CreateDepositRequest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication required."})
		return
	}

	var req models.DepositRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	response, err := ctrl.WalletSvc.InitiateDeposit(userID.(uint), req.Amount, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate deposit", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *WalletController) VerifyDeposit(c *gin.Context) {
	var req models.DepositVerifyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// In a real system, you would verify req.WebhookSignature here
	// using a shared secret with the payment gateway.
	// if !verifyWebhookSignature(req.WebhookSignature, c.Request.Body) {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
	//     return
	// }

	err := ctrl.WalletSvc.VerifyDeposit(req.PaymentGatewayTxID, req.Amount, req.Status)
	if err != nil {
		// Log the error and return a 500. Payment gateways usually retry on 5xx.
		// For user-friendly errors, you might check specific error types.
		// E.g., if errors.Is(err, service.ErrDepositNotFound) { c.JSON(http.StatusNotFound, ...) }
		if errors.Is(err, errors.New(fmt.Sprintf("deposit request not found for payment gateway transaction ID: %s", req.PaymentGatewayTxID))) { // Example, better with custom error types
			c.JSON(http.StatusNotFound, gin.H{"error": "Deposit request not found or already processed", "details": err.Error()})
		} else if errors.Is(err, errors.New(fmt.Sprintf("amount mismatch for deposit"))) { // Example
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount mismatch during verification", "details": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify deposit", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit verification processed successfully"})
}

func (ctrl *WalletController) CreateWithdrawalRequest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication required."})
		return
	}

	var req models.WithdrawalRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	response, err := ctrl.WalletSvc.InitiateWithdrawal(userID.(uint), req.Amount, req.Currency, req.BeneficiaryAccount)
	if err != nil {
		// Better error handling for insufficient funds vs. other errors
		if errors.Is(err, errors.New("insufficient funds or wallet not found")) { // Example: use custom error types for this
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate withdrawal", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *WalletController) ListWalletTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication required."})
		return
	}

	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters", "details": err.Error()})
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 { // Enforce a max limit
		params.Limit = 100
	}

	transactions, total, err := ctrl.WalletSvc.ListTransactions(userID.(uint), params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallet transactions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         params.Page,
		Limit:        params.Limit,
	})
}
