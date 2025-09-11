package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	PermSvc service.IPermissionService
	RoleSvc service.IRoleService
}

func NewPermissionController(permSvc service.IPermissionService, roleSvc service.IRoleService) *PermissionController {
	return &PermissionController{PermSvc: permSvc, RoleSvc: roleSvc}
}

func (ctrl *PermissionController) ShowAssignPage(c *gin.Context) {
	c.HTML(http.StatusOK, "assign_permissions.html", nil)
}

func (ctrl *PermissionController) GetAllPermissions(c *gin.Context) {
	groupedPerms, err := ctrl.PermSvc.GetAllGrouped()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permissions"})
		return
	}
	c.JSON(http.StatusOK, groupedPerms)
}

func (ctrl *PermissionController) GetPermissionsForRole(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if roleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := ctrl.RoleSvc.GetRoleWithPermissions(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	permissionIDs := make([]uint, len(role.Permissions))
	for i, p := range role.Permissions {
		permissionIDs[i] = p.ID
	}
	c.JSON(http.StatusOK, gin.H{"permission_ids": permissionIDs})

}

func (ctrl *PermissionController) AssignPermissionsToRole(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var payload struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.RoleSvc.AssignPermissionsToRole(uint(roleID), payload.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions updated successfully"})
}
