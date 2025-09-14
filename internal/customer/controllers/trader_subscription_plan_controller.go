package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type AdminTraderSubscriptionPlanController struct {
	service service.AdminTraderSubscriptionPlanService
}

func NewAdminTraderSubscriptionPlanController(svc service.AdminTraderSubscriptionPlanService) *AdminTraderSubscriptionPlanController {
	return &AdminTraderSubscriptionPlanController{service: svc}
}

func (c *AdminTraderSubscriptionPlanController) CreateTraderSubscriptionPlan(ctx *gin.Context) {
	var plan models.TraderSubscriptionPlan
	if err := ctx.ShouldBindJSON(&plan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPlan, err := c.service.CreateTraderSubscriptionPlan(&plan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdPlan)
}

func (c *AdminTraderSubscriptionPlanController) GetTraderSubscriptionPlanByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := c.service.GetTraderSubscriptionPlanByID(uint(id))
	if err != nil {
		if err.Error() == "trader subscription plan not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, plan)
}

// @Router /admin/trader-subscription-plans [get]
func (c *AdminTraderSubscriptionPlanController) ListTraderSubscriptionPlans(ctx *gin.Context) {
	plans, err := c.service.ListTraderSubscriptionPlans()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, plans)
}

// UpdateTraderSubscriptionPlan godoc
// @Summary Update a trader subscription plan
// @Description Admin can update details of an existing trader subscription plan.
// @Tags Admin - Trader Plans
// @Accept json
// @Produce json
// @Security BearerAuth (Admin role required)
// @Param id path int true "Plan ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} models.TraderSubscriptionPlan
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/trader-subscription-plans/{id} [put]
func (c *AdminTraderSubscriptionPlanController) UpdateTraderSubscriptionPlan(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlan, err := c.service.UpdateTraderSubscriptionPlan(uint(id), updates)
	if err != nil {
		if err.Error() == "trader subscription plan not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedPlan)
}

// DeleteTraderSubscriptionPlan godoc
// @Summary Delete a trader subscription plan
// @Description Admin can delete a trader subscription plan. (Consider soft delete instead in production)
// @Tags Admin - Trader Plans
// @Produce json
// @Security BearerAuth (Admin role required)
// @Param id path int true "Plan ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/trader-subscription-plans/{id} [delete]
func (c *AdminTraderSubscriptionPlanController) DeleteTraderSubscriptionPlan(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	err = c.service.DeleteTraderSubscriptionPlan(uint(id))
	if err != nil {
		if err.Error() == "trader subscription plan not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// ToggleTraderSubscriptionPlanStatus godoc
// @Summary Activate or deactivate a trader subscription plan
// @Description Admin can change the active status of a trader subscription plan.
// @Tags Admin - Trader Plans

func (c *AdminTraderSubscriptionPlanController) ToggleTraderSubscriptionPlanStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlan, err := c.service.ToggleTraderSubscriptionPlanStatus(uint(id), req.IsActive)
	if err != nil {
		if err.Error() == "trader subscription plan not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedPlan)
}

// Helper struct for Swagger documentation for TogglePlanStatusRequest
type TogglePlanStatusRequest struct {
	IsActive bool `json:"is_active" example:"true"`
}
