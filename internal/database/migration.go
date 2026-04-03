package database

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {

	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},

		&models.CustomerProfile{},
		&models.TraderProfile{},
		&models.Wallet{},

		&models.WalletTransaction{},
		&models.DepositRequest{},
		&models.WithdrawRequest{},
		&models.WithdrawalRequest{},

		&models.TraderSignalSubscriptionPlan{},
		&models.CustomerTraderSignalSubscription{},
		&models.AdminTraderSubscriptionPlan{},
		&models.CustomerToTraderSub{},
		&models.UserSubscription{},

		&models.MarketData{},
		&models.MarketDataAPIResponse{},
		&models.Signal{},
		&models.Trade{},
		&models.LiveTrade{},
		&models.TradeLog{},
		&models.CopySession{},

		&models.KYCDocument{},
		&models.UserKYCStatus{},

		&models.TraderPerformance{},
		&models.CommissionSetting{},

		&models.WebConfiguration{},
	)

	if err != nil {
		return err
	}

	return nil
}
