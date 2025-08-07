package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var traderAnalyticsService = service.TraderAnalyticsService{
	Repo: repository.TraderAnalyticsRepository{},
}

func GetTraderStats(c *gin.Context) {
	stats, err := traderAnalyticsService.GetAllStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trader stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GetTopRankedTraders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	stats, err := traderAnalyticsService.GetRanked(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top traders"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
