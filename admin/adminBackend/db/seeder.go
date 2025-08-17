package db

import (
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/gorm"
)

func CreateAdminSeeder(db *gorm.DB) {
	var adminCount int64
	cfg := config.AppConfig

	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)

	if adminCount == 0 {
		log.Println("No admin user found. Seeding a new admin...")

		adminEmail := cfg.AdminEmail
		adminPassword := cfg.AdminPassword

		if adminEmail == "" || adminPassword == "" {
			log.Fatal("FATAL: ADMIN_EMAIL and ADMIN_PASSWORD must be set in .env to seed the first admin.")
		}

		newAdmin := models.User{
			Name:     "Default Admin",
			Email:    adminEmail,
			Password: adminPassword,
			Role:     models.RoleAdmin,
		}

		if err := db.Create(&newAdmin).Error; err != nil {
			log.Fatalf("FATAL: Failed to seed admin user: %v", err)
		}

		log.Printf("Successfully seeded admin user with email: %s\n", adminEmail)
	} else {
		log.Println("Admin user already exists. Seeder is skipping.")
	}
}
