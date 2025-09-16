package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/gin-gonic/gin"
)

type DashboardController struct{ DashboardSvc service.IDashboardService }

func NewDashboardController(dashboardSvc service.IDashboardService) *DashboardController {
	return &DashboardController{DashboardSvc: dashboardSvc}
}

func (ctrl *DashboardController) ShowDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":        "Admin Dashboard",
		"ActiveTab":    "dashboard",
		"ActiveSubTab": "",
	})
}

func (ctrl *DashboardController) GetDashboardStats(c *gin.Context) {
	stats, err := ctrl.DashboardSvc.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (ctrl *DashboardController) GetChartData(c *gin.Context) {
	charts, err := ctrl.DashboardSvc.GetChartData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chart data"})
		return
	}
	c.JSON(http.StatusOK, charts)
}

func (ctrl *DashboardController) GetTopTraders(c *gin.Context) {
	traders, err := ctrl.DashboardSvc.GetTopTraders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top traders"})
		return
	}
	c.JSON(http.StatusOK, traders)
}

func (ctrl *DashboardController) GetLatestSignups(c *gin.Context) {
	users, err := ctrl.DashboardSvc.GetLatestSignups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest signups"})
		return
	}
	c.JSON(http.StatusOK, users)
}
