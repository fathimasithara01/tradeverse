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
		// Core Identity and Permissions
		&models.User{},       // Users must exist first
		&models.Role{},       // If Role is a separate table, it usually defines roles for users
		&models.Permission{}, // If permissions are separate from roles

		// User-related profiles and wallets
		&models.CustomerProfile{}, // Depends on User
		&models.TraderProfile{},   // Depends on User
		&models.Wallet{},          // Depends on User

		// Financial Transactions (depend on User/Wallet)
		&models.WalletTransaction{},
		&models.DepositRequest{},  // Might depend on Wallet/User
		&models.WithdrawRequest{}, // Might depend on Wallet/User

		// Trader/Customer Specific Models (Crucial ordering here)
		&models.TraderSubscriptionPlan{},     // A trader's plan. Depends on User (for TraderID).
		&models.CustomerTraderSubscription{}, // A customer's subscription to a plan. Depends on User (for CustomerID) and TraderSubscriptionPlan.
		// Generic Subscriptions (if different from CustomerTraderSubscription)
		&models.SubscriptionPlan{}, // Generic plan
		// &models.Subscription{},     // Generic subscription to a plan. Depends on User and SubscriptionPlan.
		&models.UserSubscription{}, // *Clarify: Is this distinct from Subscription? Could be redundant.*

		// Trading Related
		&models.MarketData{},
		&models.MarketDataAPIResponse{},
		&models.Signal{},      // Depends on User (TraderID)
		&models.Trade{},       // Depends on User (TraderID, potentially CustomerID)
		&models.LiveTrade{},   // Depends on Trade
		&models.TradeLog{},    // Depends on Trade
		&models.CopySession{}, // Depends on User (TraderID, CustomerID)

		// KYC
		&models.KYCDocument{},   // Depends on User
		&models.UserKYCStatus{}, // Depends on User

		// Performance, Notifications, Referrals, Admin Actions
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
