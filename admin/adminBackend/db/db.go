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

func ConnectDB() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" Failed to connect to DB: %v", err)
	}

	log.Println("Migrating parent tables (users, permissions)...")
	err = DB.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
		&models.CustomerProfile{},
		&models.TraderProfile{},
		&models.CopySession{},
		&models.TradeLog{},
		&models.Order{},
		&models.OrderItem{},
		&models.Address{},
		&models.Product{},
		&models.Trade{},
	)
	if err != nil {
		log.Printf("DATABASE MIGRATION ERROR (Step 1): %v", err)
		log.Fatal("FATAL: Failed to migrate parent tables.")
	}
	log.Println("Parent tables migrated successfully.")

	log.Println(" Connected to PostgreSQL")
}
