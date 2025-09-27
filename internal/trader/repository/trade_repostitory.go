package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TradeRepository interface {
	CreateTrade(ctx context.Context, trade *models.Trade) error
	GetTradeByID(ctx context.Context, tradeID uint) (*models.Trade, error)

	GetTradesByTraderID(ctx context.Context, traderID uint, limit, offset int) ([]models.Trade, int64, error)
	UpdateTrade(ctx context.Context, trade *models.Trade) error
	DeleteTrade(ctx context.Context, tradeID uint) error
	CloseTrade(ctx context.Context, tradeID uint, closePrice float64) (*models.Trade, error)
	CancelTrade(ctx context.Context, tradeID uint) (*models.Trade, error)
	CountTradesByTraderID(ctx context.Context, traderID uint) (int64, error)
}

type gormTradeRepository struct {
	db *gorm.DB
}

func NewGormTradeRepository(db *gorm.DB) TradeRepository {
	return &gormTradeRepository{db: db}
}

func (r *gormTradeRepository) CreateTrade(ctx context.Context, trade *models.Trade) error {
	return r.db.WithContext(ctx).Create(trade).Error
}

func (r *gormTradeRepository) GetTradeByID(ctx context.Context, tradeID uint) (*models.Trade, error) {
	var trade models.Trade
	err := r.db.WithContext(ctx).First(&trade, tradeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &trade, nil
}

func (r *gormTradeRepository) GetTradesByTraderID(ctx context.Context, traderID uint, limit, offset int) ([]models.Trade, int64, error) {
	var trades []models.Trade
	var total int64

	query := r.db.WithContext(ctx).Where("trader_id = ?", traderID)

	if err := query.Model(&models.Trade{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&trades).Error
	if err != nil {
		return nil, 0, err
	}
	return trades, total, nil
}

func (r *gormTradeRepository) UpdateTrade(ctx context.Context, trade *models.Trade) error {
	return r.db.WithContext(ctx).Save(trade).Error
}

func (r *gormTradeRepository) DeleteTrade(ctx context.Context, tradeID uint) error {
	return r.db.WithContext(ctx).Delete(&models.Trade{}, tradeID).Error
}

func (r *gormTradeRepository) CloseTrade(ctx context.Context, tradeID uint, closePrice float64) (*models.Trade, error) {
	trade, err := r.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.Status != models.TradeStatusOpen && trade.Status != models.TradeStatusPending {
		return nil, errors.New("cannot close a non-open or non-pending trade")
	}

	trade.ClosePrice = &closePrice
	trade.ClosedAt = models.TimePtr(time.Now())
	trade.Status = models.TradeStatusClosed

	pnl := (closePrice - trade.EntryPrice) * trade.Quantity * float64(trade.Leverage)
	if trade.Side == models.TradeSideSell {
		pnl = (trade.EntryPrice - closePrice) * trade.Quantity * float64(trade.Leverage)
	}
	trade.Pnl = &pnl

	err = r.db.WithContext(ctx).Save(trade).Error
	if err != nil {
		return nil, err
	}
	return trade, nil
}

func (r *gormTradeRepository) CancelTrade(ctx context.Context, tradeID uint) (*models.Trade, error) {
	trade, err := r.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.Status != models.TradeStatusPending && trade.Status != models.TradeStatusOpen {
		return nil, errors.New("only pending or open trades can be cancelled")
	}

	trade.Status = models.TradeStatusCancelled
	err = r.db.WithContext(ctx).Save(trade).Error
	if err != nil {
		return nil, err
	}
	return trade, nil
}

func (r *gormTradeRepository) CountTradesByTraderID(ctx context.Context, traderID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Trade{}).Where("trader_id = ?", traderID).Count(&count).Error
	return count, err
}
