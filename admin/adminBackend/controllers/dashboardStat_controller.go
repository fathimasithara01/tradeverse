package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	DashboardSvc *service.DashboardService
}

func NewDashboardController(dashboardSvc *service.DashboardService) *DashboardController {
	return &DashboardController{DashboardSvc: dashboardSvc}
}

func (ctrl *DashboardController) ShowDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", nil)
}

func (ctrl *DashboardController) GetDashboardStats(c *gin.Context) {
	stats, err := ctrl.DashboardSvc.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
