package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
	"github.com/fathimasithara01/tradeverse/pkg/service"
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
	user.Password = c.PostForm("Password")

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
	user.Name, user.Email, user.Password = c.PostForm("Name"), c.PostForm("Email"), c.PostForm("Password")
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
	user.Password = c.PostForm("Password")

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
	userToUpdate.Password = c.PostForm("Password")

	if userToUpdate.Role == models.RoleCustomer {
		userToUpdate.CustomerProfile.PhoneNumber = c.PostForm("PhoneNumber")
		userToUpdate.CustomerProfile.ShippingAddress = c.PostForm("ShippingAddress")
	} else if userToUpdate.Role == models.RoleTrader {
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
	if userToUpdate.Role == models.RoleTrader {
		c.Redirect(http.StatusFound, "/admin/users/all")
	} else {
		c.Redirect(http.StatusFound, "/admin/users/all")
	}
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
