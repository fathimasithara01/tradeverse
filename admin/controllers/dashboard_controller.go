package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var dashboardService = service.DashboardService{
	Repo: repository.DashboardRepository{},
}

func GetAdminDashboard(c *gin.Context) {
	stats := dashboardService.GetDashboardStats()
	c.JSON(http.StatusOK, stats)
}
