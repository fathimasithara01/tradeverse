package service

import (
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
)

type PermissionService struct {
	Repo *repository.PermissionRepository
}

func NewPermissionService(repo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{Repo: repo}
}

func (s *PermissionService) GetAllGrouped() (map[string][]models.Permission, error) {
	permissions, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.Permission)
	for _, p := range permissions {
		grouped[p.Category] = append(grouped[p.Category], p)
	}
	return grouped, nil
}
