package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// TradeRepository defines methods for interacting with trade data
type TradeRepository interface {
	CreateTrade(trade *models.Trade) error
	GetTradeByID(tradeID uint, traderID uint) (*models.Trade, error)
	ListTrades(traderID uint, pagination *models.PaginationParams) ([]models.Trade, int64, error)
	UpdateTrade(trade *models.Trade) error
	DeleteTrade(tradeID uint, traderID uint) error                // Soft delete is preferred
	FindOpenTradesByTrader(traderID uint) ([]models.Trade, error) // New: For copy trading logic, to check open positions
}

// tradeRepository implements TradeRepository using GORM
type tradeRepository struct {
	db *gorm.DB
}

// NewTradeRepository creates a new instance of TradeRepository
func NewTradeRepository(db *gorm.DB) TradeRepository {
	return &tradeRepository{db: db}
}

// CreateTrade inserts a new trade into the database
func (r *tradeRepository) CreateTrade(trade *models.Trade) error {
	return r.db.Create(trade).Error
}

// GetTradeByID retrieves a single trade by its ID and ensures it belongs to the given trader
func (r *tradeRepository) GetTradeByID(tradeID uint, traderID uint) (*models.Trade, error) {
	var trade models.Trade
	err := r.db.Where("id = ? AND trader_id = ?", tradeID, traderID).First(&trade).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Or return a custom error like ErrTradeNotFound
	}
	return &trade, err
}

// ListTrades retrieves a list of trades for a specific trader with pagination
func (r *tradeRepository) ListTrades(traderID uint, pagination *models.PaginationParams) ([]models.Trade, int64, error) {
	var trades []models.Trade
	var total int64

	query := r.db.Where("trader_id = ?", traderID)

	// Count total records
	if err := query.Model(&models.Trade{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if pagination.Limit > 0 {
		query = query.Limit(pagination.Limit).Offset((pagination.Page - 1) * pagination.Limit)
	}

	// Order by creation date descending
	err := query.Order("created_at DESC").Find(&trades).Error
	return trades, total, err
}

// UpdateTrade updates an existing trade in the database
func (r *tradeRepository) UpdateTrade(trade *models.Trade) error {
	return r.db.Save(trade).Error // Save updates all fields. Use r.db.Model(trade).Updates(...) for partial updates
}

// DeleteTrade soft deletes a trade (sets DeletedAt timestamp)
func (r *tradeRepository) DeleteTrade(tradeID uint, traderID uint) error {
	// Ensure the trade belongs to the trader before deleting
	return r.db.Where("id = ? AND trader_id = ?", tradeID, traderID).Delete(&models.Trade{}).Error
}

// FindOpenTradesByTrader retrieves all open trades for a specific trader
func (r *tradeRepository) FindOpenTradesByTrader(traderID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.Where("trader_id = ? AND status = ?", traderID, models.TradeStatusOpen).Find(&trades).Error
	return trades, err
}
