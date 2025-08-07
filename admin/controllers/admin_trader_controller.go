package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var traderService = service.TraderService{
	Repo: repository.TraderRepository{},
}

func GetAllTraders(c *gin.Context) {
	traders, err := traderService.GetAllTraders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traders"})
		return
	}
	c.JSON(http.StatusOK, traders)
}

func ToggleBanTrader(c *gin.Context) {
	idParam := c.Param("id")
	traderID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID"})
		return
	}

	err = traderService.ToggleBan(uint(traderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trader status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trader ban status updated successfully"})
}
