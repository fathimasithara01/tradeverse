package service

import (
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
)

type RoleService struct {
	Repo           *repository.RoleRepository
	PermissionRepo *repository.PermissionRepository
}

func NewRoleService(repo *repository.RoleRepository, permRepo *repository.PermissionRepository) *RoleService {
	return &RoleService{Repo: repo, PermissionRepo: permRepo}
}

func (s *RoleService) CreateRole(role *models.Role, loggedInUserID uint) error {
	// role.CreatedByID = loggedInUserID
	return s.Repo.Create(role)
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.Repo.FindAll()
}

func (s *RoleService) GetRoleByID(id uint) (models.Role, error) {
	return s.Repo.FindByID(id)
}

func (s *RoleService) UpdateRole(role *models.Role) error {
	return s.Repo.Update(role)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.Repo.Delete(id)
}

func (s *RoleService) GetRoleWithPermissions(id uint) (models.Role, error) {
	return s.Repo.FindByIDWithPermissions(id)
}

func (s *RoleService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	role, err := s.Repo.FindByID(roleID)
	if err != nil {
		return err
	}

	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		if err := s.PermissionRepo.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}
	}

	return s.Repo.UpdatePermissions(&role, permissions)
}
