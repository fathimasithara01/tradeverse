package db

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
		return nil, err
	}
	log.Println("Parent tables migrated successfully.")

	return db, err
}
