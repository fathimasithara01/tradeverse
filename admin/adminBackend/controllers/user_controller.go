package controllers

import (
	"net/http"
	"strings"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserSvc      *service.UserService
	DashboardSvc *service.DashboardService
}

func NewUserController(userSvc *service.UserService, dashboardSvc *service.DashboardService) *UserController {
	return &UserController{
		UserSvc:      userSvc,
		DashboardSvc: dashboardSvc,
	}
}

func (ctrl *UserController) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func (ctrl *UserController) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (ctrl *UserController) LoginAdmin(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	// The login service now returns the token and the user object
	token, user, err := ctrl.UserSvc.Login(email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	// **IMPORTANT**: Check if the logged-in user is actually an admin
	if user.Role != models.RoleAdmin {
		c.HTML(http.StatusForbidden, "login.html", gin.H{"error": "You do not have permission to access the admin panel."})
		return
	}

	c.SetCookie("admin_token", token, 3600*24, "/", config.AppConfig.CookieDomain, false, true)
	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func (ctrl *UserController) ShowDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", nil)
}

func (ctrl *UserController) GetDashboardStats(c *gin.Context) {
	stats, err := ctrl.DashboardSvc.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dashboard stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
