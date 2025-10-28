package controllers

import (
	"log"
	"net/http"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserSvc service.IUserService
}

func NewAuthController(userSvc service.IUserService) *AuthController {
	return &AuthController{UserSvc: userSvc}
}

func (ctrl *AuthController) ShowLoginPage(c *gin.Context) {
	errorMessage := c.Query("error")
	c.HTML(http.StatusOK, "login.html", gin.H{"error": errorMessage})
}

func (ctrl *AuthController) ShowCustomerRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register_customer.html", nil)
}

func (ctrl *AuthController) ShowTraderRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register_trader.html", nil)
}

func (ctrl *AuthController) LoginUser(c *gin.Context) {
	email, password := c.PostForm("email"), c.PostForm("password")

	token, user, err := ctrl.UserSvc.Login(email, password)
	if err != nil {
		log.Printf("[LOGIN FAILED] for user '%s': %v", email, err)
		c.Redirect(http.StatusFound, "/login?error="+err.Error())
		return
	}

	log.Printf("[LOGIN SUCCESS] User '%s' logged in successfully. Role: %s\n", user.Email, user.Role)
	c.SetCookie("admin_token", token, 86400, "/", config.AppConfig.Cookie.Domain, true, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}








func (ctrl *AuthController) RegisterCustomer(c *gin.Context) {
	var user models.User
	var profile models.CustomerProfile
	user.Name = c.PostForm("Name")
	user.Email = c.PostForm("Email")
	rawPassword := c.PostForm("Password")
	confirmPassword := c.PostForm("ConfirmPassword")
	profile.Phone = c.PostForm("PhoneNumber")

	if rawPassword != confirmPassword {
		c.HTML(http.StatusBadRequest, "register_customer.html", gin.H{"error": "Passwords do not match."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[REGISTER FAILED] Hashing password for user '%s' failed: %v", user.Email, err)
		c.HTML(http.StatusInternalServerError, "register_customer.html", gin.H{"error": "Failed to process password."})
		return
	}
	user.Password = string(hashedPassword)

	if err := ctrl.UserSvc.RegisterCustomer(user, profile); err != nil {
		log.Printf("[REGISTER FAILED] Customer registration for '%s' failed: %v", user.Email, err)
		c.HTML(http.StatusBadRequest, "register_customer.html", gin.H{"error": err.Error()})
		return
	}
	log.Printf("[REGISTER SUCCESS] Customer '%s' registered successfully.", user.Email)
	c.Redirect(http.StatusFound, "/login?message=Registration successful! Please log in.")
}

func (ctrl *AuthController) RegisterTrader(c *gin.Context) {
	var user models.User
	var profile models.TraderProfile
	user.Name = c.PostForm("Name")
	user.Email = c.PostForm("Email")
	rawPassword := c.PostForm("Password")
	confirmPassword := c.PostForm("ConfirmPassword")
	profile.CompanyName = c.PostForm("CompanyName")

	if rawPassword != confirmPassword {
		c.HTML(http.StatusBadRequest, "register_trader.html", gin.H{"error": "Passwords do not match."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[REGISTER FAILED] Hashing password for trader '%s' failed: %v", user.Email, err)
		c.HTML(http.StatusInternalServerError, "register_trader.html", gin.H{"error": "Failed to process password."})
		return
	}
	user.Password = string(hashedPassword)

	if err := ctrl.UserSvc.RegisterTrader(user, profile); err != nil {
		log.Printf("[REGISTER FAILED] Trader registration for '%s' failed: %v", user.Email, err)
		c.HTML(http.StatusBadRequest, "register_trader.html", gin.H{"error": err.Error()})
		return
	}
	log.Printf("[REGISTER SUCCESS] Trader '%s' registered successfully. Awaiting admin approval.", user.Email)
	c.Redirect(http.StatusFound, "/login?message=Trader registration successful! Awaiting admin approval.")
}
