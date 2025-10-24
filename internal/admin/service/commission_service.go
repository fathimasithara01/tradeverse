package service

import (
	"fmt"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICommissionService interface {
	SetPlatformCommissionPercentage(adminID uint, percentage float64) (*models.AdminCommissionResponsePayload, error)
	GetPlatformCommissionPercentage() (*models.AdminCommissionResponsePayload, error)
}

type CommissionService struct {
	CommissionRepo repository.ICommissionRepository
	DB             *gorm.DB // In case transactions are needed for more complex operations
}

// NewCommissionService creates a new instance of CommissionService.
func NewCommissionService(commissionRepo repository.ICommissionRepository, db *gorm.DB) *CommissionService {
	return &CommissionService{
		CommissionRepo: commissionRepo,
		DB:             db,
	}
}

// SetPlatformCommissionPercentage sets or updates the global platform commission percentage.
func (s *CommissionService) SetPlatformCommissionPercentage(adminID uint, percentage float64) (*models.AdminCommissionResponsePayload, error) {
	if percentage < 0 || percentage > 100 {
		return nil, fmt.Errorf("commission percentage must be between 0 and 100")
	}

	commissionSetting := &models.CommissionSetting{
		Key:         "trader_subscription_commission_percentage",
		Value:       percentage,
		Description: "Percentage of subscription fees taken as platform commission from traders.",
		UpdatedBy:   adminID,
	}

	err := s.CommissionRepo.CreateOrUpdateCommissionSetting(commissionSetting)
	if err != nil {
		return nil, fmt.Errorf("failed to set platform commission percentage: %w", err)
	}

	// Retrieve the updated setting to get its ID and LastUpdated time
	updatedSetting, err := s.CommissionRepo.GetCommissionSettingByKey("trader_subscription_commission_percentage")
	if err != nil {
		// This should not happen immediately after update/create but as a fallback
		return nil, fmt.Errorf("failed to retrieve updated commission setting: %w", err)
	}

	return &models.AdminCommissionResponsePayload{
		ID:                   updatedSetting.ID,
		CommissionPercentage: updatedSetting.Value,
		LastUpdated:          updatedSetting.LastUpdated,
		UpdatedBy:            updatedSetting.UpdatedBy,
		Description:          updatedSetting.Description,
	}, nil
}

func (s *CommissionService) GetPlatformCommissionPercentage() (*models.AdminCommissionResponsePayload, error) {
	setting, err := s.CommissionRepo.GetCommissionSettingByKey("trader_subscription_commission_percentage")
	if err != nil {
		if err.Error() == "commission setting not found" {
			return &models.AdminCommissionResponsePayload{
				CommissionPercentage: 0.0, // Default to 0 if not set
				Description:          "Commission percentage not explicitly set, defaulting to 0.",
			}, nil
		}
		return nil, fmt.Errorf("failed to retrieve platform commission percentage: %w", err)
	}

	return &models.AdminCommissionResponsePayload{
		ID:                   setting.ID,
		CommissionPercentage: setting.Value,
		LastUpdated:          setting.LastUpdated,
		UpdatedBy:            setting.UpdatedBy,
		Description:          setting.Description,
	}, nil
}
