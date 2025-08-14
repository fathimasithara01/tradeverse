package controllers

import (
	"log"
	"net/http"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserSvc *service.UserService
}

func NewAuthController(userSvc *service.UserService) *AuthController {
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
func (ctrl *AuthController) ShowAdminRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register_admin.html", nil)
}

func (ctrl *AuthController) LoginUser(c *gin.Context) {
	email, password := c.PostForm("email"), c.PostForm("password")
	token, user, err := ctrl.UserSvc.Login(email, password)
	if err != nil {
		log.Printf("[LOGIN FAILED] Attempt for user '%s' failed: %s\n", email, err.Error())
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("[LOGIN SUCCESS] User '%s' logged in. Generating token for UserID: %d with Role: %s\n", user.Email, user.ID, user.Role)

	if user.Role == models.RoleAdmin {
		c.SetCookie("admin_token", token, 3600*24, "/", config.AppConfig.CookieDomain, false, true)
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.SetCookie("user_token", token, 3600*24, "/", config.AppConfig.CookieDomain, false, true)
		c.Redirect(http.StatusFound, "/dashboard")
	}
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

func (ctrl *AuthController) RegisterAdmin(c *gin.Context) {
	var user models.User
	user.Name, user.Email, user.Password = c.PostForm("Name"), c.PostForm("Email"), c.PostForm("Password")
	if c.PostForm("Password") != c.PostForm("ConfirmPassword") {
		c.HTML(http.StatusBadRequest, "register_admin.html", gin.H{"error": "Passwords do not match."})
		return
	}
	if _, err := ctrl.UserSvc.RegisterAdmin(user); err != nil {
		c.HTML(http.StatusBadRequest, "register_admin.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/login")
}
