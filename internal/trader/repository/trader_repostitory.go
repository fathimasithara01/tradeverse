package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TradeRepository interface {
	GetAllTrades(traderID uint, limit, offset int) ([]models.Trade, int64, error)
	GetTradeByID(id uint, traderID uint) (*models.Trade, error)
	CreateTrade(ctx context.Context, req models.TradeRequest) (models.Trade, error)
	UpdateTrade(trade *models.Trade) error
	DeleteTrade(id uint, traderID uint) error
}

type tradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) TradeRepository {
	return &tradeRepository{db: db}
}

func (r *tradeRepository) GetAllTrades(traderID uint, limit, offset int) ([]models.Trade, int64, error) {
	var trades []models.Trade
	var count int64

	err := r.db.Model(&models.Trade{}).
		Where("trader_id = ?", traderID).
		Count(&count).
		Limit(limit).Offset(offset).
		Order("created_at desc").
		Find(&trades).Error
	return trades, count, err
}

func (r *tradeRepository) GetTradeByID(id uint, traderID uint) (*models.Trade, error) {
	var trade models.Trade
	err := r.db.Where("id = ? AND trader_id = ?", id, traderID).First(&trade).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &trade, err
}

func (r *tradeRepository) CreateTrade(ctx context.Context, req models.TradeRequest) (models.Trade, error) {
	trade := models.Trade{
		TraderID:        req.TraderID,
		Symbol:          req.Symbol,
		TradeType:       models.TradeType(req.TradeType), // string → TradeType
		Side:            models.TradeSide(req.Side),      // string → TradeSide
		EntryPrice:      req.EntryPrice,
		Quantity:        req.Quantity,
		Leverage:        uint(req.Leverage),            // int → uint
		StopLossPrice:   floatPtr(req.StopLossPrice),   // float64 → *float64
		TakeProfitPrice: floatPtr(req.TakeProfitPrice), // float64 → *float64
		Status:          "OPEN",
		OpenedAt:        timePtr(time.Now()), // time.Time → *time.Time
	}

	if err := r.db.WithContext(ctx).Create(&trade).Error; err != nil {
		return models.Trade{}, err
	}
	return trade, nil
}

func (r *tradeRepository) UpdateTrade(trade *models.Trade) error {
	return r.db.Save(trade).Error
}

func (r *tradeRepository) DeleteTrade(id uint, traderID uint) error {
	return r.db.Where("id = ? AND trader_id = ?", id, traderID).
		Delete(&models.Trade{}).Error
}

func floatPtr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}
