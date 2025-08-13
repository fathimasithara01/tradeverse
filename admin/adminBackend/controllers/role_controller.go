package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	RoleSvc *service.RoleService
}

func NewRoleController(roleSvc *service.RoleService) *RoleController {
	return &RoleController{RoleSvc: roleSvc}
}

func (ctrl *RoleController) ShowRolesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage_roles.html", nil)
}

func (ctrl *RoleController) ShowAddRolePage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_role.html", nil)
}

func (ctrl *RoleController) ShowEditRolePage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	role, err := ctrl.RoleSvc.GetRoleByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Role not found"})
		return
	}
	c.HTML(http.StatusOK, "edit_role.html", gin.H{"Role": role})
}

func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	role.Name = c.PostForm("Name")

	loggedInUserID, exists := c.Get("userID")
	if !exists {
		log.Println("[ERROR] 'userID' not found in Gin context.")
		c.HTML(http.StatusUnauthorized, "add_role.html", gin.H{"error": "User ID could not be determined from token."})
		return
	}

	userID, ok := loggedInUserID.(uint)
	if !ok || userID == 0 {
		log.Printf("[ERROR] 'userID' in context has wrong type or is zero. Value: %v", loggedInUserID)
		c.HTML(http.StatusInternalServerError, "add_role.html", gin.H{"error": "User ID has an invalid format."})
		return
	}

	log.Printf("[INFO] Attempting to create role with CreatedByID: %d\n", userID)
	if err := ctrl.RoleSvc.CreateRole(&role, userID); err != nil {
		log.Printf("[ERROR] Service failed to create role: %s\n", err.Error())
		c.HTML(http.StatusBadRequest, "add_role.html", gin.H{"error": "Failed to create role: " + err.Error()})
		return
	}

	log.Println("[SUCCESS] Role created successfully.")
	c.Redirect(http.StatusFound, "/admin/roles")
}

func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	role, err := ctrl.RoleSvc.GetRoleByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "edit_role.html", gin.H{"error": "Role not found"})
		return
	}

	role.Name = c.PostForm("Name")
	if err := ctrl.RoleSvc.UpdateRole(&role); err != nil {
		c.HTML(http.StatusInternalServerError, "edit_role.html", gin.H{"error": "Failed to update role", "Role": role})
		return
	}
	c.Redirect(http.StatusFound, "/admin/roles")
}

func (ctrl *RoleController) GetRoles(c *gin.Context) {
	roles, err := ctrl.RoleSvc.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve roles"})
		return
	}
	if roles == nil {
		roles = make([]models.Role, 0)
	}
	c.JSON(http.StatusOK, roles)
}

func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID."})
		return
	}
	if err := ctrl.RoleSvc.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
