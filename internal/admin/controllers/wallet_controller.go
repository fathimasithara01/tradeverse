package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
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

type PaginationParams struct {
	Page        int    `form:"page,default=1"`
	Limit       int    `form:"limit,default=10"`
	SearchQuery string `form:"search"`
}

// ✅ 1. HTML Page Renderer
func (ctrl *AdminWalletController) ShowAllCustomerTransactionsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "all_customer_transactions.html", gin.H{
		"Title":        "All Customer Transactions",
		"ActiveTab":    "financials",
		"ActiveSubTab": "all_transactions",
	})
	fmt.Println("✅ Finished rendering all_customer_transactions.html")
}

func (ctrl *AdminWalletController) AdminGetAllCustomerTransactions(c *gin.Context) {
	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		fmt.Printf("❌ Error binding query params: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid pagination parameters",
			"details": err.Error(),
		})
		return
	}

	// Default values
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	fmt.Printf("📄 Fetching transactions | Page: %d | Limit: %d | Search: %s\n",
		pagination.Page, pagination.Limit, pagination.Search)

	transactions, total, err := ctrl.AdminWalletService.GetAllCustomerTransactionsWithUserDetails(pagination)
	if err != nil {
		fmt.Printf("❌ Error retrieving customer transactions: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve all customer transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         pagination.Page,
		"limit":        pagination.Limit,
	})
}

func (ctrl *AdminWalletController) ShowAdminWalletPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_wallet.html", gin.H{
		"Title":        "Admin Wallet",
		"ActiveTab":    "financials",
		"ActiveSubTab": "admin_wallet",
	})
	fmt.Println("Finished rendering admin_wallet.html (if no error occurred)")
}
func (ctrl *AdminWalletController) ShowAdminWalletTransactionPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_wallet_transactions.html", gin.H{
		"Title":        "Admin Wallet Transaction",
		"ActiveTab":    "financials",
		"ActiveSubTab": "admin_wallet_transactions",
	})
	fmt.Println("Finished rendering admin_wallet.html (if no error occurred)")
}

func (ctrl *AdminWalletController) AdminGetAllPlatformTransactions(c *gin.Context) {
	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, total, err := ctrl.AdminWalletService.GetAllWalletTransactions(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all platform wallet transactions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
	})
}
func (ctrl *AdminWalletController) GetAdminWalletSummary(c *gin.Context) {
	summary, err := ctrl.AdminWalletService.GetAdminWalletSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin wallet summary", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

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

func (ctrl *AdminWalletController) AdminVerifyDeposit(c *gin.Context) {
	var req models.DepositVerifyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func (ctrl *AdminWalletController) AdminGetWalletTransactions(c *gin.Context) {
	var pagination models.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, total, err := ctrl.AdminWalletService.AdminGetWalletTransactions(pagination) // Changed this line
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
