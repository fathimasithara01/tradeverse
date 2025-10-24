package controllers

import (
	"fmt"
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type CommissionController struct {
	CommissionService service.ICommissionService
}

func NewCommissionController(commissionService service.ICommissionService) *CommissionController {
	return &CommissionController{
		CommissionService: commissionService,
	}
}

func (ctrl *CommissionController) ShowCommissionSettingsPage(c *gin.Context) {
	currentSetting, err := ctrl.CommissionService.GetPlatformCommissionPercentage()
	var currentPercentage float64
	if err != nil {
		fmt.Printf("Error fetching current commission setting for page: %v\n", err)
		currentPercentage = 0.0
	} else {
		currentPercentage = currentSetting.CommissionPercentage
	}

	c.HTML(http.StatusOK, "admin_commission_settings.html", gin.H{
		"Title":             "Commission Settings",
		"ActiveTab":         "financials",
		"ActiveSubTab":      "commission",
		"CurrentCommission": currentPercentage,
	})
}

// GetCommissionSettings handles fetching the current commission percentage via API.
func (ctrl *CommissionController) GetCommissionSettings(c *gin.Context) {
	setting, err := ctrl.CommissionService.GetPlatformCommissionPercentage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve commission settings", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, setting)
}

// UpdateCommissionSettings handles updating the platform commission percentage via API.
func (ctrl *CommissionController) UpdateCommissionSettings(c *gin.Context) {
	var payload models.AdminCommissionRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.MustGet("userID").(uint) // Assuming your auth middleware sets userID

	response, err := ctrl.CommissionService.SetPlatformCommissionPercentage(adminID, payload.CommissionPercentage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update commission settings", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commission percentage updated successfully", "data": response})
}
