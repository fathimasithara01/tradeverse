package controllers

import (
	"log"
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserSvc service.IUserService
}

func NewAuthController(userSvc service.IUserService) *AuthController {
	return &AuthController{UserSvc: userSvc}
}

func (ctrl *AuthController) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
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
		log.Printf("[LOGIN FAILED] Invalid credentials for user '%s'", email)
		c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
		return
	}

	log.Printf("[LOGIN SUCCESS] User '%s' logged in successfully. Role: %s\n", user.Email, user.Role)
	c.SetCookie("admin_token", token, 86400, "/", config.AppConfig.CookieDomain, false, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (ctrl *AuthController) RegisterCustomer(c *gin.Context) {
	var user models.User
	var profile models.CustomerProfile
	user.Name, user.Email, user.Password = c.PostForm("Name"), c.PostForm("Email"), c.PostForm("Password")
	profile.PhoneNumber = c.PostForm("PhoneNumber")
	if c.PostForm("Password") != c.PostForm("ConfirmPassword") {
		c.HTML(http.StatusBadRequest, "register_customer.html", gin.H{"error": "Passwords do not match."})
		return
	}
	if err := ctrl.UserSvc.RegisterCustomer(user, profile); err != nil {
		c.HTML(http.StatusBadRequest, "register_customer.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/login")
}

func (ctrl *AuthController) RegisterTrader(c *gin.Context) {
	var user models.User
	var profile models.TraderProfile
	user.Name, user.Email, user.Password = c.PostForm("Name"), c.PostForm("Email"), c.PostForm("Password")
	profile.CompanyName = c.PostForm("CompanyName")
	if c.PostForm("Password") != c.PostForm("ConfirmPassword") {
		c.HTML(http.StatusBadRequest, "register_trader.html", gin.H{"error": "Passwords do not match."})
		return
	}
	if err := ctrl.UserSvc.RegisterTrader(user, profile); err != nil {
		c.HTML(http.StatusBadRequest, "register_trader.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/login")
}
