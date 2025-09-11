package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IPermissionRepository interface {
	FindAll() ([]models.Permission, error)
}

type PermissionRepository struct{ DB *gorm.DB }

func NewPermissionRepository(db *gorm.DB) IPermissionRepository {
	return &PermissionRepository{DB: db}
}

func (r *PermissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.DB.Order("category asc, name asc").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
