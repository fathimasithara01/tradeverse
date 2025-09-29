package service

import (
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TraderService struct {
	traderRepo *customerrepo.TraderRepository
	db         *gorm.DB
}

func NewTraderService(traderRepo *customerrepo.TraderRepository, db *gorm.DB) *TraderService {
	return &TraderService{
		traderRepo: traderRepo,
		db:         db,
	}
}

func (s *TraderService) ListApprovedTraders(filters map[string]interface{}, sortBy string, sortOrder string, page, pageSize int) ([]models.TraderProfile, int64, error) {
	traders, total, err := s.traderRepo.FindApprovedTradersWithUsers(filters, sortBy, sortOrder, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return traders, total, nil
}

func (s *TraderService) GetTraderDetails(traderID uint) (*models.TraderProfile, error) {
	traderProfile, err := s.traderRepo.FindTraderProfileWithUser(traderID)
	if err != nil {
		return nil, err
	}
	return traderProfile, nil
}

func (s *TraderService) GetTraderPerformance(traderID uint) (interface{}, error) {

	return map[string]string{"message": "Trader performance history not yet implemented"}, nil
}
