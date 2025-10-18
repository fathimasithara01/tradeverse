package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type TraderProfileController struct {
	traderService service.ITraderProfileService
}

func NewTraderProfileController(traderService service.ITraderProfileService) *TraderProfileController {
	return &TraderProfileController{traderService: traderService}
}

func (c *TraderProfileController) GetTraderProfile(ctx *gin.Context) {

	userID := ctx.MustGet("userID").(uint)

	profile, err := c.traderService.GetProfile(userID)
	if err != nil {
		if err == service.ErrTraderProfileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

type CreateTraderProfileRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	CompanyName string `json:"company_name" binding:"required,max=100"`
	Bio         string `json:"bio" binding:"required"`
}

func (c *TraderProfileController) CreateTraderProfile(ctx *gin.Context) {
	var req CreateTraderProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userID := ctx.MustGet("userID").(uint)

	profile, err := c.traderService.CreateProfile(userID, req.Name, req.CompanyName, req.Bio)
	if err != nil {
		if err == service.ErrTraderProfileExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrPermissionDenied {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, profile)
}

type UpdateTraderProfileRequest struct {
	Name        *string `json:"name,omitempty" binding:"max=100"`
	CompanyName *string `json:"company_name,omitempty" binding:"max=100"`
	Bio         *string `json:"bio,omitempty"`
}

func (c *TraderProfileController) UpdateTraderProfile(ctx *gin.Context) {
	var req UpdateTraderProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userID := ctx.MustGet("userID").(uint)

	existingProfile, err := c.traderService.GetProfile(userID)
	if err != nil {
		if err == service.ErrTraderProfileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	profileID := existingProfile.ID

	updatedProfile, err := c.traderService.UpdateProfile(userID, profileID, req.Name, req.CompanyName, req.Bio)
	if err != nil {
		if err == service.ErrTraderProfileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorized {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedProfile)
}

func (c *TraderProfileController) DeleteTraderProfile(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)

	existingProfile, err := c.traderService.GetProfile(userID)
	if err != nil {
		if err == service.ErrTraderProfileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	profileID := existingProfile.ID

	err = c.traderService.DeleteProfile(userID, profileID)
	if err != nil {
		if err == service.ErrTraderProfileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorized {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *TraderProfileController) ApproveTraderProfile(ctx *gin.Context) {

	profileIDStr := ctx.Param("id")
	profileID, err := strconv.ParseUint(profileIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	profile, err := c.traderService.GetProfile(uint(profileID))

	profile, err = c.traderService.GetProfile(uint(profileID))
	if err != nil {
		if err == service.ErrTraderProfileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profile.Status = models.StatusApproved
	profile, err = c.traderService.UpdateProfile(profile.UserID, profile.ID, &profile.Name, &profile.CompanyName, &profile.Bio)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile status: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}
