package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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

	if role.Name == "" {
		c.HTML(http.StatusBadRequest, "add_role.html", gin.H{"error": "Role name cannot be empty."})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Authentication error: User ID not found."})
		return
	}

	loggedInUserID, ok := userIDVal.(uint)
	if !ok || loggedInUserID == 0 {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Authentication error: Invalid User ID."})
		return
	}

	if err := ctrl.RoleSvc.CreateRole(&role, loggedInUserID); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == "23505" {
			if strings.Contains(pgErr.Message, "uni_roles_name") {
				c.HTML(http.StatusBadRequest, "add_role.html", gin.H{
					"error": "Failed to create role: A role with this name already exists.",
					"Name":  role.Name,
				})
				return
			}
		}

		c.HTML(http.StatusInternalServerError, "add_role.html", gin.H{
			"error": "An unexpected error occurred. Please try again.",
			"Name":  role.Name,
		})
		return
	}

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

	c.JSON(http.StatusOK, roles)
}

func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.RoleSvc.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
