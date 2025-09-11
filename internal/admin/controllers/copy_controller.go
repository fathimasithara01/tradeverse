package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/gin-gonic/gin"
)

type CopyController struct{ CopySvc service.ICopyService }

func NewCopyController(copySvc service.ICopyService) *CopyController {
	return &CopyController{CopySvc: copySvc}
}

func (ctrl *CopyController) GetCopyStatus(c *gin.Context) {
	followerID, _ := c.Get("userID")
	masterID, _ := strconv.ParseUint(c.Param("masterID"), 10, 32)

	isActive, _ := ctrl.CopySvc.GetCopyStatus(followerID.(uint), uint(masterID))
	c.JSON(http.StatusOK, gin.H{"is_copying": isActive})
}

func (ctrl *CopyController) StartCopying(c *gin.Context) {
	followerID, _ := c.Get("userID")
	var payload struct {
		MasterID uint `json:"master_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master ID"})
		return
	}

	err := ctrl.CopySvc.StartCopying(followerID.(uint), payload.MasterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start copying session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Now copying trader!"})
}

func (ctrl *CopyController) StopCopying(c *gin.Context) {
	followerID, _ := c.Get("userID")
	var payload struct {
		MasterID uint `json:"master_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master ID"})
		return
	}

	err := ctrl.CopySvc.StopCopying(followerID.(uint), payload.MasterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop copying session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stopped copying trader."})
}
