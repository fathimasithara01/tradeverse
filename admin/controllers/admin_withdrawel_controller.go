package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var withdrawalService = service.WithdrawalService{
	Repo: repository.WithdrawalRepository{},
}

func GetPendingWithdrawals(c *gin.Context) {
	list, err := withdrawalService.GetPending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch pending withdrawals"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func ApproveWithdrawal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := withdrawalService.Approve(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Approval failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal approved"})
}

func RejectWithdrawal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var body struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note is required"})
		return
	}

	err := withdrawalService.Reject(uint(id), body.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rejection failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal rejected"})
}
