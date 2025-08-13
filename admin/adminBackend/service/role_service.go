package service

import (
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
)

type RoleService struct {
	Repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{Repo: repo}
}

func (s *RoleService) CreateRole(role *models.Role, loggedInUserID uint) error {
	role.CreatedByID = loggedInUserID
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
