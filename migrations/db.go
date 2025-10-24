package migrations

import (
	"fmt"
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)

	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
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
	)
	if err != nil {
		log.Printf("DATABASE MIGRATION ERROR (Step 1): %v", err)
		return nil, err
	}
	log.Println("Parent tables migrated successfully.")

	return db, err
}
