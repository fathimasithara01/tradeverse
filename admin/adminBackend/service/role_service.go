package service

import (
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
)

type RoleService struct {
	RoleRepo *repository.RoleRepository
	PermRepo *repository.PermissionRepository
	UserRepo *repository.UserRepository // Add this line

}

func NewRoleService(roleRepo *repository.RoleRepository, permRepo *repository.PermissionRepository, userRepo *repository.UserRepository) *RoleService {
	return &RoleService{RoleRepo: roleRepo, PermRepo: permRepo, UserRepo: userRepo}

}

type RoleResponse struct {
	ID            uint   `json:"ID"`
	Name          string `json:"name"`
	CreatedByID   uint   `json:"created_by_id"`
	CreatedByName string `json:"createdByName"` // The new field for the creator's name
}

func (s *RoleService) GetAllRoles() ([]RoleResponse, error) {
	roles, err := s.RoleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return []RoleResponse{}, nil
	}

	userIDs := make([]uint, 0)
	for _, role := range roles {
		if role.CreatedByID > 0 {
			userIDs = append(userIDs, role.CreatedByID)
		}
	}

	users, err := s.UserRepo.FindByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[uint]string)
	for _, user := range users {
		userMap[user.ID] = user.Name
	}

	var responses []RoleResponse
	for _, role := range roles {
		responses = append(responses, RoleResponse{
			ID:            role.ID,
			Name:          role.Name,
			CreatedByID:   role.CreatedByID,
			CreatedByName: userMap[role.CreatedByID],
		})
	}
	return responses, nil
}
func (s *RoleService) CreateRole(role *models.Role, loggedInUserID uint) error {
	role.CreatedByID = loggedInUserID
	return s.RoleRepo.Create(role)
}

// func (s *RoleService) GetAllRoles() ([]models.Role, error) {
// 	return s.RoleRepo.FindAll()
// }

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

func (s *RoleService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	role, err := s.RoleRepo.FindByID(roleID)
	if err != nil {
		return err
	}
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		if err := s.PermRepo.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}
	}
	return s.RoleRepo.UpdatePermissions(&role, permissions)
}
