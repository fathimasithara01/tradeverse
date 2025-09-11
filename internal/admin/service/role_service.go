package service

import (
	"fmt"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type IRoleService interface {
	GetAllRoles() ([]models.Role, error)
	GetAllRolesWithUsers() ([]RoleWithUsers, error)
	CreateRole(role *models.Role, loggedInUserID uint) error
	GetRoleByID(id uint) (models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error
	GetRoleWithPermissions(id uint) (models.Role, error)
	AssignPermissionsToRole(roleID uint, permissionIDs []uint) error
	RoleHasPermission(roleID uint, permissionName string) (bool, error)
}

type RoleService struct {
	RoleRepo       repository.IRoleRepository
	PermissionRepo repository.IPermissionRepository
	UserRepo       repository.IUserRepository
}

func NewRoleService(
	roleRepo repository.IRoleRepository,
	permissionRepo repository.IPermissionRepository,
	userRepo repository.IUserRepository,
) IRoleService {
	return &RoleService{
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
		UserRepo:       userRepo,
	}
}

type RoleWithUsers struct {
	models.Role
	Users []models.User `json:"users"`
}

func (s *RoleService) GetAllRolesWithUsers() ([]RoleWithUsers, error) {
	roles, err := s.RoleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return []RoleWithUsers{}, nil
	}

	var rolesWithUsers []RoleWithUsers

	for _, role := range roles {
		users, err := s.UserRepo.FindByRole(models.UserRole(role.Name))
		if err != nil {
			fmt.Printf("Warning: could not fetch users for role %s: %v\n", role.Name, err)
			users = []models.User{}
		}
		rolesWithUsers = append(rolesWithUsers, RoleWithUsers{
			Role:  role,
			Users: users,
		})
	}
	return rolesWithUsers, nil
}

func (s *RoleService) CreateRole(role *models.Role, loggedInUserID uint) error {
	role.CreatedByID = loggedInUserID
	return s.RoleRepo.Create(role)
}

func (s *RoleService) GetRoleByID(id uint) (models.Role, error) {
	return s.RoleRepo.FindByID(id)
}

func (s *RoleService) UpdateRole(role *models.Role) error {
	return s.RoleRepo.Update(role)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.RoleRepo.Delete(id)
}

func (s *RoleService) GetRoleWithPermissions(id uint) (models.Role, error) {
	return s.RoleRepo.FindByIDWithPermissions(id)
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.RoleRepo.FindAll()
}
func (s *RoleService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	role, err := s.RoleRepo.FindByID(roleID)
	if err != nil {
		return err
	}

	_ = s.PermissionRepo
	var permissions []models.Permission

	return s.RoleRepo.UpdatePermissions(&role, permissions)
}

func (s *RoleService) RoleHasPermission(roleID uint, permissionName string) (bool, error) {
	return s.RoleRepo.RoleHasPermission(roleID, permissionName)
}
