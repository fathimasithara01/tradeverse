package customerrepo

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TraderRepository struct {
	db *gorm.DB
}

func NewTraderRepository(db *gorm.DB) *TraderRepository {
	return &TraderRepository{db: db}
}

func (r *TraderRepository) FindApprovedTradersWithUsers(filters map[string]interface{}, sortBy string, sortOrder string, page, pageSize int) ([]models.TraderProfile, int64, error) {
	var traderProfiles []models.TraderProfile
	var total int64

	query := r.db.Preload(clause.Associations).Where("status = ?", models.StatusApproved)

	if companyName, ok := filters["company_name"]; ok && companyName != "" {
		query = query.Where("LOWER(company_name) LIKE LOWER(?)", "%"+companyName.(string)+"%")
	}
	if isVerified, ok := filters["is_verified"]; ok {
		query = query.Where("is_verified = ?", isVerified)
	}

	query.Model(&models.TraderProfile{}).Count(&total)

	if sortBy != "" {
		order := sortBy
		if sortOrder == "desc" {
			order = sortBy + " DESC"
		} else {
			order = sortBy + " ASC"
		}
		query = query.Order(order)
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Limit(pageSize).Offset(offset)
	}

	err := query.Find(&traderProfiles).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range traderProfiles {
		var user models.User
		if err := r.db.Where("id = ?", traderProfiles[i].UserID).First(&user).Error; err != nil {
			continue
		}
	}

	return traderProfiles, total, nil
}

func (r *TraderRepository) FindTraderProfileWithUser(traderID uint) (*models.TraderProfile, error) {
	var traderProfile models.TraderProfile
	err := r.db.Preload(clause.Associations).
		Where("id = ? AND status = ?", traderID, models.StatusApproved).
		First(&traderProfile).Error
	if err != nil {
		return nil, err
	}
	return &traderProfile, nil
}
