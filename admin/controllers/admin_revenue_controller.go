package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var revenueService = service.RevenueService{
	Repo: repository.RevenueRepository{},
}

func GetMonthlyRevenue(c *gin.Context) {
	data, err := revenueService.GetMonthly()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch revenue data"})
		return
	}
	c.JSON(http.StatusOK, data)
}
