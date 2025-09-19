package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TraderController struct {
	traderService *service.TraderService
}

func NewTraderController(traderService *service.TraderService) *TraderController {
	return &TraderController{
		traderService: traderService,
	}
}

func (ctrl *TraderController) ListTraders(c *gin.Context) {
	filters := make(map[string]interface{})
	if companyName := c.Query("company_name"); companyName != "" {
		filters["company_name"] = companyName
	}
	if isVerifiedStr := c.Query("is_verified"); isVerifiedStr != "" {
		isVerified, err := strconv.ParseBool(isVerifiedStr)
		if err == nil {
			filters["is_verified"] = isVerified
		}
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	traders, total, err := ctrl.traderService.ListApprovedTraders(filters, sortBy, sortOrder, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve traders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      traders,
		"total":     total,
		"page":      page,
		"last_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (ctrl *TraderController) GetTraderDetails(c *gin.Context) {
	traderIDStr := c.Param("id")
	traderID, err := strconv.ParseUint(traderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID"})
		return
	}

	trader, err := ctrl.traderService.GetTraderDetails(uint(traderID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trader not found or not approved"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve trader details"})
		return
	}

	c.JSON(http.StatusOK, trader)
}

func (ctrl *TraderController) GetTraderPerformance(c *gin.Context) {
	traderIDStr := c.Param("id")
	traderID, err := strconv.ParseUint(traderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID"})
		return
	}

	_, err = ctrl.traderService.GetTraderDetails(uint(traderID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trader not found or not approved"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify trader existence"})
		return
	}

	performance, err := ctrl.traderService.GetTraderPerformance(uint(traderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve trader performance"})
		return
	}

	c.JSON(http.StatusOK, performance)
}
