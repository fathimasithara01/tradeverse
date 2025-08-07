package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/gin-gonic/gin"
)

func GetAllRevenueSplits(c *gin.Context) {
	var splits []models.RevenueSplit
	if err := db.DB.Find(&splits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch revenue splits"})
		return
	}
	c.JSON(http.StatusOK, splits)
}
