package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type LiveTradeRepository interface {
	Create(liveTrade *models.LiveTrade) error
	GetActive(traderID uint) ([]models.LiveTrade, error)
}

type liveTradeRepo struct {
	db *gorm.DB
}

func NewLiveTradeRepository(db *gorm.DB) LiveTradeRepository {
	return &liveTradeRepo{db: db}
}

func (r *liveTradeRepo) Create(liveTrade *models.LiveTrade) error {
	return r.db.Create(liveTrade).Error
}

func (r *liveTradeRepo) GetActive(traderID uint) ([]models.LiveTrade, error) {
	var trades []models.LiveTrade
	err := r.db.Where("trader_id = ? AND status = ?", traderID, "OPEN").Find(&trades).Error
	return trades, err
}
