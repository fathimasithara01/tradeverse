package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var traderStatsService = service.TraderStatsService{
	Repo: repository.TraderStatsRepository{},
}

func GetTraderRankings(c *gin.Context) {
	rankings, err := traderStatsService.GetAllTraderRankings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trader rankings"})
		return
	}
	c.JSON(http.StatusOK, rankings)
}

func GetTraderBadge(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	badge, err := traderStatsService.GetBadgeForTrader(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get badge"})
		return
	}
	c.JSON(http.StatusOK, badge)
}
