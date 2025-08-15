package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserSvc *service.UserService
}

func NewUserController(userSvc *service.UserService) *UserController {
	return &UserController{UserSvc: userSvc}
}

func (ctrl *UserController) ShowUsersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_users.html", nil)
}

func (ctrl *UserController) ShowAddUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_user.html", nil)
}

func (ctrl *UserController) ShowEditUserPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.UserSvc.GetUserByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "edit_user.html", gin.H{"User": user})
}
func (ctrl *UserController) CreateCustomer(c *gin.Context) {
	var user models.User
	var profile models.CustomerProfile
	user.Name, user.Email, user.Password = c.PostForm("Name"), c.PostForm("Email"), c.PostForm("Password")
	profile.PhoneNumber = c.PostForm("PhoneNumber")

	if err := ctrl.UserSvc.RegisterCustomer(user, profile); err != nil {
		c.HTML(http.StatusBadRequest, "add_user.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid user ID"})
		return
	}

	userToUpdate, err := ctrl.UserSvc.GetUserByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}

	userToUpdate.Name = c.PostForm("Name")
	userToUpdate.Email = c.PostForm("Email")
	userToUpdate.Password = c.PostForm("Password") // The service will handle hashing if not empty.

	if userToUpdate.Role == models.RoleCustomer {
		userToUpdate.CustomerProfile.PhoneNumber = c.PostForm("PhoneNumber")
		userToUpdate.CustomerProfile.ShippingAddress = c.PostForm("ShippingAddress")
	} else if userToUpdate.Role == models.RoleTrader {
		userToUpdate.TraderProfile.CompanyName = c.PostForm("CompanyName")
		userToUpdate.TraderProfile.Bio = c.PostForm("Bio")
	}

	// Call the service to perform the update.
	if err := ctrl.UserSvc.UpdateUser(&userToUpdate); err != nil {
		c.HTML(http.StatusInternalServerError, "edit_user.html", gin.H{
			"error": "Failed to update user.",
			"User":  userToUpdate, // Send the user back to re-populate the form
		})
		return
	}

	// Redirect based on the user's role to the correct management page.
	if userToUpdate.Role == models.RoleTrader {
		c.Redirect(http.StatusFound, "/admin/users/traders")
	} else {
		c.Redirect(http.StatusFound, "/admin/users/customers")
	}
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.UserSvc.GetUsersByRole(models.RoleCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	if users == nil {
		users = make([]models.User, 0)
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := ctrl.UserSvc.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
