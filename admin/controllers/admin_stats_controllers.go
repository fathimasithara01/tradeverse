package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var statsService = service.StatsService{
	Repo: repository.StatsRepository{},
}

func GetPlanStats(c *gin.Context) {
	stats, err := statsService.GetPlanStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plan usage stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
