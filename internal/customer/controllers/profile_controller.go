package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	UserSvc service.IUserService
}

func NewProfileController(userSvc service.IUserService) *ProfileController {
	return &ProfileController{UserSvc: userSvc}
}

func (ctrl *ProfileController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	user, err := ctrl.UserSvc.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

type UpdateProfileRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"omitempty,email"`
	PhoneNumber *string `json:"phone_number"`
}

func (ctrl *ProfileController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := ctrl.UserSvc.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updatedUser := models.User{}
	updatedProfile := models.CustomerProfile{}

	if req.Name != nil {
		updatedUser.Name = *req.Name
	}
	if req.Email != nil {
		updatedUser.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		updatedProfile.Phone = *req.PhoneNumber
	}

	err = ctrl.UserSvc.UpdateCustomerProfile(existingUser.ID, updatedUser, updatedProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (ctrl *ProfileController) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	err := ctrl.UserSvc.DeleteUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
