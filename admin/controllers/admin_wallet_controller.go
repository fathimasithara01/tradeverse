package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var walletService = service.WalletService{
	Repo: repository.WalletRepository{},
}

func GetWalletDetails(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("user_id"))
	wallet, txs, err := walletService.GetWallet(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wallet": wallet, "transactions": txs})
}

func CreditWallet(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("user_id"))

	var body struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := walletService.Credit(uint(uid), body.Amount, body.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to credit wallet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Wallet credited"})
}

func DebitWallet(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("user_id"))

	var body struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := walletService.Debit(uint(uid), body.Amount, body.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to debit wallet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Wallet debited"})
}
