package controllers

import (
	"html/template"
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	Service *service.AdminService
}

func NewAdminController(db *gorm.DB) *AdminController {
	repo := repository.NewAdminRepository(db)
	svc := service.NewAdminService(repo)
	return &AdminController{Service: svc}
}

func (ctrl *AdminController) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func (ctrl *AdminController) RegisterAdmin(c *gin.Context) {
	admin := models.Admin{
		Name:     c.PostForm("name"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	_, err := ctrl.Service.Register(admin)
	if err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Registration failed"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func (ctrl *AdminController) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (ctrl *AdminController) LoginAdmin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, err := ctrl.Service.Login(email, password)
	if err != nil || token == "" {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid login"})
		return
	}

	c.SetCookie("admin_token", token, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func (ctrl *AdminController) LogoutAdmin(c *gin.Context) {
	c.SetCookie("admin_token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func (ctrl *AdminController) AdminDashboard(c *gin.Context) {
	tmpl, _ := template.ParseFiles("templates/dashboard.html")
	tmpl.Execute(c.Writer, gin.H{"title": "Admin Dashboard"})
}
