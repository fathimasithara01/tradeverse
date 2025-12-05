// internal/admin/repository/web_configuration_repository.go
package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type WebConfigurationRepository interface {
	GetWebConfiguration() (*models.WebConfiguration, error)
	UpdateWebConfiguration(config *models.WebConfiguration) error
}

type webConfigurationRepository struct {
	db *gorm.DB
}

func NewWebConfigurationRepository(db *gorm.DB) WebConfigurationRepository {
	return &webConfigurationRepository{db: db}
}

func (r *webConfigurationRepository) GetWebConfiguration() (*models.WebConfiguration, error) {
	var config models.WebConfiguration
	if err := r.db.First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *webConfigurationRepository) UpdateWebConfiguration(config *models.WebConfiguration) error {
	return r.db.Save(config).Error
}
