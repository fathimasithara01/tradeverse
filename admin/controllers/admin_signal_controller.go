package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var signalService = service.SignalService{
	Repo: repository.SignalRepository{},
}

func GetAllSignals(c *gin.Context) {
	signals, err := signalService.GetAllSignals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve signals"})
		return
	}
	c.JSON(http.StatusOK, signals)
}

func DeactivateSignal(c *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signal ID"})
		return
	}

	err := signalService.Deactivate(uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate signal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signal deactivated successfully"})
}

func GetPendingSignals(c *gin.Context) {
	signals, err := signalService.GetPending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending signals"})
		return
	}
	c.JSON(http.StatusOK, signals)
}

func ApproveSignal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := signalService.Approve(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve signal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Signal approved"})
}

func RejectSignal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := signalService.Reject(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject signal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Signal rejected"})
}
