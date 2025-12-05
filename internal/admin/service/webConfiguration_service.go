package service

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type WebConfigurationService interface {
	GetWebConfiguration() (*models.WebConfiguration, error)
	UpdateWebConfiguration(primaryCountry, primaryCurrency, primaryTimezone string) error
}

type webConfigurationService struct {
	webConfigRepo repository.WebConfigurationRepository
}

func NewWebConfigurationService(webConfigRepo repository.WebConfigurationRepository) WebConfigurationService {
	return &webConfigurationService{webConfigRepo: webConfigRepo}
}

func (s *webConfigurationService) GetWebConfiguration() (*models.WebConfiguration, error) {
	return s.webConfigRepo.GetWebConfiguration()
}

func (s *webConfigurationService) UpdateWebConfiguration(primaryCountry, primaryCurrency, primaryTimezone string) error {
	config, err := s.webConfigRepo.GetWebConfiguration()
	if err != nil {
		return err
	}

	config.PrimaryCountry = primaryCountry
	config.PrimaryCurrency = primaryCurrency
	config.PrimaryTimezone = primaryTimezone

	return s.webConfigRepo.UpdateWebConfiguration(config)
}
