package service

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
)

type IPermissionService interface {
	GetAllGrouped() (map[string][]models.Permission, error)
}

type PermissionService struct {
	Repo repository.IPermissionRepository
}

func NewPermissionService(repo repository.IPermissionRepository) IPermissionService {
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
