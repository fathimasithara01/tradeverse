package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type KYCController struct {
	KYCSvc service.KYCServicer
}

func NewKYCController(kycSvc service.KYCServicer) *KYCController {
	return &KYCController{KYCSvc: kycSvc}
}

func (ctrl *KYCController) SubmitKYCDocuments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication required."})
		return
	}

	var req models.SubmitKYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := ctrl.KYCSvc.SubmitKYCDocuments(userID.(uint), req.DocumentType, req.DocumentURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit KYC documents", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC documents submitted successfully. Verification status updated to PENDING."})
}

func (ctrl *KYCController) GetKYCStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication required."})
		return
	}

	statusResponse, err := ctrl.KYCSvc.GetKYCStatus(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve KYC status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statusResponse)
}
