package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/models"
	"github.com/gin-gonic/gin"
)

func (ctrl *UserController) ShowUsersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_users.html", nil)
}

func (ctrl *UserController) ShowAddUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_user.html", nil)
}

func (ctrl *UserController) ShowEditUserPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
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

// --- API and Form Handlers ---

func (ctrl *UserController) CreateCustomer(c *gin.Context) {
	var user models.User
	var profile models.CustomerProfile

	user.Name = c.PostForm("Name")
	user.Email = c.PostForm("Email")
	user.Password = c.PostForm("Password")
	profile.ShippingAddress = c.PostForm("ShippingAddress")
	profile.PhoneNumber = c.PostForm("PhoneNumber")

	if err := ctrl.UserSvc.CreateCustomer(user, profile); err != nil {
		c.HTML(http.StatusBadRequest, "add_user.html", gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	// Fetch the existing user to update.
	userToUpdate, err := ctrl.UserSvc.GetUserByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}

	// Update fields from form data.
	userToUpdate.Name = c.PostForm("Name")
	userToUpdate.Email = c.PostForm("Email")
	userToUpdate.Password = c.PostForm("Password") // Service will handle if it's empty.
	userToUpdate.CustomerProfile.ShippingAddress = c.PostForm("ShippingAddress")
	userToUpdate.CustomerProfile.PhoneNumber = c.PostForm("PhoneNumber")

	if err := ctrl.UserSvc.UpdateUser(&userToUpdate); err != nil {
		c.HTML(http.StatusInternalServerError, "edit_user.html", gin.H{
			"error": "Failed to update user.",
			"User":  userToUpdate,
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.UserSvc.GetUsersByRole(models.RoleCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteUser is the API for the delete button.
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
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
