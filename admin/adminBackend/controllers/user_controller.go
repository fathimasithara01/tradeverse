package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/models"
	"github.com/gin-gonic/gin"
)

func (ctrl *UserController) ShowCustomersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_customers.html", nil)
}

func (ctrl *UserController) ShowTradersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_traders.html", nil)
}

func (ctrl *UserController) ShowAddTraderPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_trader.html", nil)
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

func (ctrl *UserController) ShowTraderApprovalPage(c *gin.Context) {
	c.HTML(http.StatusOK, "trader_approval.html", nil)
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
