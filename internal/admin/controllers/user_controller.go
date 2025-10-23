package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type UserController struct{ UserSvc service.IUserService }

func NewUserController(userSvc service.IUserService) *UserController {
	return &UserController{UserSvc: userSvc}
}

func (ctrl *UserController) ShowUsersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_users.html", nil)
}
func (ctrl *UserController) ShowCustomersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_customers.html", nil)
}
func (ctrl *UserController) ShowTradersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_traders.html", nil)
}
func (ctrl *UserController) ShowAddCustomerPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_customer.html", nil)
}
func (ctrl *UserController) ShowAddTraderPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_trader.html", nil)
}
func (ctrl *UserController) ShowAddInternalUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_internal_user.html", nil)

}
func (ctrl *UserController) ShowTraderApprovalPage(c *gin.Context) {
	c.HTML(http.StatusOK, "trader_approval.html", nil)
}
func (ctrl *UserController) ShowAssignRolePage(c *gin.Context) {
	c.HTML(http.StatusOK, "assign_role.html", nil)
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

func (ctrl *UserController) CreateTrader(c *gin.Context) {
	var user models.User
	var profile models.TraderProfile

	user.Name = c.PostForm("Name")
	user.Email = c.PostForm("Email")
	rawPassword := c.PostForm("Password")

	if err := user.SetPassword(rawPassword); err != nil {
		c.HTML(http.StatusInternalServerError, "add_trader.html", gin.H{"error": "Failed to process password."})
		return
	}

	profile.CompanyName = c.PostForm("CompanyName")
	profile.Bio = c.PostForm("Bio")

	if err := ctrl.UserSvc.CreateTraderByAdmin(user, profile); err != nil {
		c.HTML(http.StatusBadRequest, "add_trader.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users/traders")
}

func (ctrl *UserController) CreateCustomer(c *gin.Context) {
	var user models.User
	var profile models.CustomerProfile
	user.Name, user.Email = c.PostForm("Name"), c.PostForm("Email")
	rawPassword := c.PostForm("Password")

	if err := user.SetPassword(rawPassword); err != nil {
		c.HTML(http.StatusInternalServerError, "add_customer.html", gin.H{"error": "Failed to process password."})
		return
	}

	profile.PhoneNumber = c.PostForm("PhoneNumber")

	if err := ctrl.UserSvc.RegisterCustomer(user, profile); err != nil {
		c.HTML(http.StatusBadRequest, "add_customer.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users/customers")
}

func (ctrl *UserController) CreateInternalUser(c *gin.Context) {
	var user models.User
	user.Name = c.PostForm("Name")
	user.Email = c.PostForm("Email")
	rawPassword := c.PostForm("Password")

	if err := user.SetPassword(rawPassword); err != nil {
		c.HTML(http.StatusInternalServerError, "add_internal_user.html", gin.H{"error": "Failed to process password."})
		return
	}

	if _, err := ctrl.UserSvc.CreateInternalUser(user); err != nil {
		c.HTML(http.StatusBadRequest, "add_internal_user.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users/all")
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
	newPassword := c.PostForm("Password")

	if newPassword != "" {
		if err := userToUpdate.SetPassword(newPassword); err != nil {
			c.HTML(http.StatusInternalServerError, "edit_user.html", gin.H{
				"error": "Failed to process new password.",
				"User":  userToUpdate,
			})
			return
		}
	}

	if userToUpdate.Role == models.RoleCustomer {
		if userToUpdate.CustomerProfile.ID == 0 {

		}
		userToUpdate.CustomerProfile.Name = c.PostForm("Name")
		userToUpdate.CustomerProfile.PhoneNumber = c.PostForm("PhoneNumber")
		userToUpdate.CustomerProfile.ShippingAddress = c.PostForm("ShippingAddress")
	} else if userToUpdate.Role == models.RoleTrader {
		if userToUpdate.TraderProfile.ID == 0 {
		}
		userToUpdate.TraderProfile.CompanyName = c.PostForm("CompanyName")
		userToUpdate.TraderProfile.Bio = c.PostForm("Bio")
	}

	if err := ctrl.UserSvc.UpdateUser(&userToUpdate); err != nil {
		c.HTML(http.StatusInternalServerError, "edit_user.html", gin.H{
			"error": "Failed to update user.",
			"User":  userToUpdate,
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users/all")
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

func (ctrl *UserController) GetCustomers(c *gin.Context) {
	customers, err := ctrl.UserSvc.GetUsersByRole(models.RoleCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}
	if customers == nil {
		customers = make([]models.User, 0)
	}
	c.JSON(http.StatusOK, customers)
}

func (ctrl *UserController) GetTraders(c *gin.Context) {
	traders, err := ctrl.UserSvc.GetUsersByRole(models.RoleTrader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve traders"})
		return
	}
	if traders == nil {
		traders = make([]models.User, 0)
	}
	c.JSON(http.StatusOK, traders)
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	allUsers, err := ctrl.UserSvc.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all users"})
		return
	}

	if allUsers == nil {
		allUsers = make([]models.User, 0)
	}

	c.JSON(http.StatusOK, allUsers)
}

func (ctrl *UserController) GetAllUsersAdvanced(c *gin.Context) {
	var options repository.UserQueryOptions

	if err := c.ShouldBindQuery(&options); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	paginatedResult, err := ctrl.UserSvc.GetAllUsersAdvanced(options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, paginatedResult)
}

func (ctrl *UserController) GetPendingTraders(c *gin.Context) {
	traders, err := ctrl.UserSvc.GetTradersByStatus(models.StatusPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pending traders"})
		return
	}
	if traders == nil {
		traders = make([]models.User, 0)
	}
	c.JSON(http.StatusOK, traders)
}

func (ctrl *UserController) GetApprovedTraders(c *gin.Context) {
	traders, err := ctrl.UserSvc.GetTradersByStatus(models.StatusApproved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve approved traders"})
		return
	}
	if traders == nil {
		traders = make([]models.User, 0)
	}
	c.JSON(http.StatusOK, traders)
}

func (ctrl *UserController) ApproveTrader(c *gin.Context) {
	traderID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if traderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID"})
		return
	}
	if err := ctrl.UserSvc.ApproveTrader(uint(traderID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve trader"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trader approved successfully"})
}

func (ctrl *UserController) RejectTrader(c *gin.Context) {
	traderID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if traderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader ID"})
		return
	}
	if err := ctrl.UserSvc.RejectTrader(uint(traderID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject trader"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trader rejected successfully"})
}

func (ctrl *UserController) AssignRoleToUser(c *gin.Context) {
	var payload struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: User ID and Role ID are required."})
		return
	}

	if err := ctrl.UserSvc.AssignRoleToUser(payload.UserID, payload.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully."})
}

func (ctrl *UserController) GetUsersForRoleAssignment(c *gin.Context) {
	users, err := ctrl.UserSvc.GetAllUsersWithRole()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	if users == nil {
		users = make([]models.User, 0)
	}

	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) ShowAdminProfileViewPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_profile_view.html", gin.H{"Title": "Admin Profile"})
}

func (ctrl *UserController) ShowAdminProfileEditPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_profile_edit.html", gin.H{"Title": "Edit Admin Profile"})
}
func (ctrl *UserController) GetAdminProfileAPI(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	adminID, ok := userID.(uint)
	if !ok || adminID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}

	user, err := ctrl.UserSvc.GetAdminProfile(adminID)
	if err != nil {
		log.Printf("[ERROR] GetAdminProfileAPI for ID %d: %v", adminID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch admin profile: %v", err.Error())})
		return
	}

	// Ensure password is never sent in the API response
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// UpdateAdminProfileAPI handles updating the admin profile via API.
func (ctrl *UserController) UpdateAdminProfileAPI(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	adminID, ok := userID.(uint)
	if !ok || adminID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}

	var req service.AdminUpdateProfileRequest
	// For multipart/form-data (which includes file uploads), ShouldBind handles it
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request data: %v", err.Error())})
		return
	}

	if err := ctrl.UserSvc.UpdateAdminProfile(adminID, req); err != nil {
		log.Printf("[ERROR] UpdateAdminProfileAPI for ID %d: %v", adminID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update admin profile: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin profile updated successfully"})
}

// ChangeAdminPasswordAPI handles changing the admin password via API.
func (ctrl *UserController) ChangeAdminPasswordAPI(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	adminID, ok := userID.(uint)
	if !ok || adminID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.UserSvc.ChangeAdminPassword(adminID, req.OldPassword, req.NewPassword); err != nil {
		log.Printf("[ERROR] ChangeAdminPasswordAPI for ID %d: %v", adminID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // Use 401 for incorrect password, 500 for other issues
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// ShowAdminSettingsPage renders a placeholder for admin settings.
func (ctrl *UserController) ShowAdminSettingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_settings.html", gin.H{"Title": "Admin Settings"})
}
