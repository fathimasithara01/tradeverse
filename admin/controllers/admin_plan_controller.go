package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var planService = service.PlanService{
	Repo: repository.PlanRepository{},
}

func GetPendingPlans(c *gin.Context) {
	plans, err := planService.GetPendingPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func ApprovePlan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := planService.ApprovePlan(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Approval failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan approved"})
}

func RejectPlan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := planService.RejectPlan(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rejection failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan rejected"})
}

func GetAllPlans(c *gin.Context) {
	plans, err := planService.GetAllPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func CreatePlan(c *gin.Context) {
	var input models.Plan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := planService.CreatePlan(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func UpdatePlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input models.Plan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := planService.UpdatePlan(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan updated"})
}

func DeletePlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := planService.DeletePlan(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deletion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
}
