package customerrepo

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type KYCRepository interface {
	CreateKYCDocument(doc *models.KYCDocument) error
	FindUserKYCStatus(userID uint) (*models.UserKYCStatus, error)
	CreateUserKYCStatus(status *models.UserKYCStatus) error
	UpdateUserKYCStatus(status *models.UserKYCStatus) error
}

type kycRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) KYCRepository {
	return &kycRepository{db: db}
}

func (r *kycRepository) CreateKYCDocument(doc *models.KYCDocument) error {
	return r.db.Create(doc).Error
}

func (r *kycRepository) FindUserKYCStatus(userID uint) (*models.UserKYCStatus, error) {
	var status models.UserKYCStatus
	err := r.db.Where("user_id = ?", userID).First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User's KYC status not found, not an error
		}
		return nil, err
	}
	return &status, nil
}

func (r *kycRepository) CreateUserKYCStatus(status *models.UserKYCStatus) error {
	return r.db.Create(status).Error
}

func (r *kycRepository) UpdateUserKYCStatus(status *models.UserKYCStatus) error {
	// Only update specific fields
	return r.db.Model(status).Updates(map[string]interface{}{
		"status":            status.Status,
		"reason":            status.Reason,
		"last_updated_by":   status.LastUpdatedBy,
		"last_updated_date": status.LastUpdatedDate,
	}).Error
}
