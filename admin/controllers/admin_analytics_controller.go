package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var analyticsService = service.AnalyticsService{
	Repo: repository.AnalyticsRepository{},
}

func GetSignalAnalytics(c *gin.Context) {
	data := analyticsService.GetSignalStats()
	c.JSON(http.StatusOK, data)
}

func GetTraderPerformance(c *gin.Context) {
	stats, err := analyticsService.GetTraderStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
