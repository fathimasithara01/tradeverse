package db

import (
	"fmt"
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" Failed to connect to DB: %v", err)
	}

	// err = DB.AutoMigrate(&models.Admin{})
	err = DB.AutoMigrate(
		&models.User{},
		&models.CustomerProfile{}, // Add this
		&models.TraderProfile{},
		&models.Role{}, // Add this
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Address{},
		// &models.Plan{},
		// &models.Signal{},
		// &models.Subscription{},
		// &models.Payment{},
		// &models.Announcement{},
		// &models.Log{},
		// &models.Follower{},
		// &models.Withdrawal{},
		// &models.Wallet{},
		// &models.WalletTransaction{},
		// &models.Notification{},
		// &models.Payment{},
		// &models.RevenueSplit{},
	)
	if err != nil {
		log.Fatal("error ")
	}

	log.Println(" Connected to PostgreSQL")
}
