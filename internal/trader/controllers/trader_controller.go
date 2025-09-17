package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/customer/service" // Re-use customer service for now
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type TraderController interface {
	GetMyTraderProfile(c *gin.Context)
	// Add other trader-specific endpoints here (e.g., UpdateProfile, ListMyTrades)
}

type traderController struct {
	customerService service.CustomerService // Re-using for trader profile logic
	// You might have a dedicated TraderService later
}

func NewTraderController(cs service.CustomerService) TraderController {
	return &traderController{customerService: cs}
}

func (ctrl *traderController) GetMyTraderProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Assumes AuthMiddleware sets userID

	// Ensure the user is actually a trader
	user, err := ctrl.customerService.repo.GetUserByID(userID) // Access repo via service
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user details: " + err.Error()})
		return
	}
	if user == nil || user.Role != models.RoleTrader {
		c.JSON(http.StatusForbidden, gin.H{"message": "access denied: not a trader"})
		return
	}

	profile, err := ctrl.customerService.GetTraderProfileDetails(userID)
	if err != nil {
		if err.Error() == "trader profile not found for this user" { // Specific error check
			c.JSON(http.StatusNotFound, gin.H{"message": "trader profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve trader profile: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}
