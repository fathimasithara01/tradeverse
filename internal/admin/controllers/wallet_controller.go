package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models" // Assuming models has Wallet structs
	"github.com/gin-gonic/gin"
)

type AdminWalletController struct {
	AdminWalletService service.IAdminWalletService
}

func NewAdminWalletController(adminWalletService service.IAdminWalletService) *AdminWalletController {
	return &AdminWalletController{
		AdminWalletService: adminWalletService,
	}
}

func (ctrl *AdminWalletController) ShowAdminWalletPage(c *gin.Context) {
	fmt.Println("Attempting to render admin_wallet.html") // Add this line for debugging
	c.HTML(http.StatusOK, "admin_wallet.html", gin.H{
		"Title":        "Admin Wallet",
		"ActiveTab":    "financials",
		"ActiveSubTab": "admin_wallet",
	})
	fmt.Println("Finished rendering admin_wallet.html (if no error occurred)") // Add this line
}

// GetAdminWalletSummary retrieves the admin's wallet balance and details.
func (ctrl *AdminWalletController) GetAdminWalletSummary(c *gin.Context) {
	summary, err := ctrl.AdminWalletService.GetAdminWalletSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin wallet summary", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// AdminInitiateDeposit simulates an admin initiating a deposit to their wallet.
func (ctrl *AdminWalletController) AdminInitiateDeposit(c *gin.Context) {
	var req models.DepositRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctrl.AdminWalletService.AdminInitiateDeposit(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate admin deposit", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// AdminVerifyDeposit simulates a webhook callback for a deposit.
func (ctrl *AdminWalletController) AdminVerifyDeposit(c *gin.Context) {
	var req models.DepositVerifyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming the request includes the DepositID to verify
	depositIDStr := c.Param("deposit_id")
	depositID, err := strconv.ParseUint(depositIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit ID"})
		return
	}

	res, err := ctrl.AdminWalletService.AdminVerifyDeposit(uint(depositID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin deposit", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// AdminRequestWithdrawal handles an admin's request to withdraw funds.
func (ctrl *AdminWalletController) AdminRequestWithdrawal(c *gin.Context) {
	var req models.WithdrawalRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctrl.AdminWalletService.AdminRequestWithdrawal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process admin withdrawal request", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// AdminGetWalletTransactions retrieves all transactions for the admin's wallet.
func (ctrl *AdminWalletController) AdminGetWalletTransactions(c *gin.Context) {
	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// You might want to add filtering options here as well (e.g., by type, status, date range)
	transactions, total, err := ctrl.AdminWalletService.AdminGetWalletTransactions(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin wallet transactions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
	})
}

// AdminApproveOrRejectWithdrawal allows admin to approve/reject customer withdrawal requests
func (ctrl *AdminWalletController) AdminApproveOrRejectWithdrawal(c *gin.Context) {
	withdrawalIDStr := c.Param("id")
	withdrawalID, err := strconv.ParseUint(withdrawalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdrawal request ID"})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required,oneof=approve reject"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Action == "approve" {
		err = ctrl.AdminWalletService.ApproveWithdrawalRequest(uint(withdrawalID))
	} else {
		err = ctrl.AdminWalletService.RejectWithdrawalRequest(uint(withdrawalID))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to %s withdrawal request", req.Action), "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Withdrawal request %sed successfully", req.Action)})
}

// GetPendingWithdrawals retrieves all pending withdrawal requests (from customers)
func (ctrl *AdminWalletController) GetPendingWithdrawals(c *gin.Context) {
	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdrawals, total, err := ctrl.AdminWalletService.GetPendingWithdrawalRequests(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending withdrawal requests", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"withdrawals": withdrawals,
		"total":       total,
		"page":        pagination.Page,
		"limit":       pagination.Limit,
	})
}
