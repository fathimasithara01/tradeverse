package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService service.WalletService
}

func NewWalletController(walletService service.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

func (ctrl *WalletController) GetBalance(c *gin.Context) {
	userID := c.GetUint("userID")

	wallet, err := ctrl.walletService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Wallet balance retrieved", wallet)
}

func (ctrl *WalletController) Deposit(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := ctrl.walletService.Deposit(c.Request.Context(), userID, req.Amount)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Deposit successful", tx)
}

func (ctrl *WalletController) Withdraw(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := ctrl.walletService.Withdraw(c.Request.Context(), userID, req.Amount)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Withdrawal successful", tx)
}

func (ctrl *WalletController) TransactionHistory(c *gin.Context) {
	userID := c.GetUint("userID")

	txs, err := ctrl.walletService.GetTransactionHistory(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Transaction history retrieved", txs)
}
