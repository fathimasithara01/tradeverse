package repository

import (
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICommissionRepository interface {
	CreateOrUpdateCommissionSetting(setting *models.CommissionSetting) error
	GetCommissionSettingByKey(key string) (*models.CommissionSetting, error)
	GetPlatformCommissionPercentage() (float64, error)
}

type CommissionRepository struct {
	DB *gorm.DB
}

func NewCommissionRepository(db *gorm.DB) *CommissionRepository {
	return &CommissionRepository{DB: db}
}

func (r *CommissionRepository) CreateOrUpdateCommissionSetting(setting *models.CommissionSetting) error {
	var existingSetting models.CommissionSetting
	result := r.DB.Where("key = ?", setting.Key).First(&existingSetting)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return r.DB.Create(setting).Error
		}
		return fmt.Errorf("failed to check existing commission setting: %w", result.Error)
	}

	existingSetting.Value = setting.Value
	existingSetting.Description = setting.Description
	existingSetting.UpdatedBy = setting.UpdatedBy
	return r.DB.Save(&existingSetting).Error
}

func (r *CommissionRepository) GetCommissionSettingByKey(key string) (*models.CommissionSetting, error) {
	var setting models.CommissionSetting
	err := r.DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("commission setting not found")
		}
		return nil, fmt.Errorf("failed to retrieve commission setting: %w", err)
	}
	return &setting, nil
}

func (r *CommissionRepository) GetPlatformCommissionPercentage() (float64, error) {
	setting, err := r.GetCommissionSettingByKey("trader_subscription_commission_percentage")
	if err != nil {
		if errors.Is(err, errors.New("commission setting not found")) {
			return 0.0, nil
		}
		return 0.0, fmt.Errorf("failed to get platform commission percentage: %w", err)
	}
	return setting.Value, nil
}
