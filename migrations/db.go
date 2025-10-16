package migrations

import (
	"fmt"
	"log"

	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

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

		&models.TraderSubscriptionPlan{},
		&models.CustomerTraderSubscription{},
		&models.AdminSubscriptionPlan{},
		&models.Subscription{},
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

		&models.TraderPerformance{}, // Depends on User (TraderID)
		&models.Notification{},      // Depends on User
		&models.Referral{},          // Depends on User (ReferrerID, RefereeID)
		&models.AdminActionLog{},    // Depends on User (AdminID)
	)
	if err != nil {
		log.Printf("DATABASE MIGRATION ERROR (Step 1): %v", err)
		return nil, err
	}
	log.Println("Parent tables migrated successfully.")

	return db, err
}
