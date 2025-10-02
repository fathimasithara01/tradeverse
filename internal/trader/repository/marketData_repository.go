package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type MarketDataRepository interface {
	Create(data *models.MarketData) error
	GetBySymbol(symbol string) (*models.MarketData, error)
}

type marketDataRepository struct {
	db *gorm.DB
}

func NewMarketDataRepository(db *gorm.DB) MarketDataRepository {
	return &marketDataRepository{db: db}
}

func (r *marketDataRepository) Create(data *models.MarketData) error {
	return r.db.Create(data).Error
}

func (r *marketDataRepository) GetBySymbol(symbol string) (*models.MarketData, error) {
	var md models.MarketData
	err := r.db.Where("symbol = ?", symbol).First(&md).Error
	if err != nil {
		return nil, err
	}
	return &md, nil
}
