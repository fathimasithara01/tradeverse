package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/admin/service" // This import might be incorrect if you have a customer-specific user service
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserSvc service.IUserService // Assuming this is the shared user service, otherwise correct the import
}

func NewAuthController(userSvc service.IUserService) *AuthController {
	return &AuthController{UserSvc: userSvc}
}

type SignupRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number"`
}

func (ctrl *AuthController) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
		// Do NOT directly assign req.Password here. It must be hashed.
	}

	// Hash the password before passing it to the service
	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	profile := models.CustomerProfile{
		PhoneNumber: req.PhoneNumber,
	}

	if err := ctrl.UserSvc.RegisterCustomer(user, profile); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer registration successful"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// The UserSvc.Login method should handle password comparison using bcrypt
	token, user, err := ctrl.UserSvc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // Return the actual error from service
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}
