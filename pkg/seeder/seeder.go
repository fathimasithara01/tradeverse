package seeder

import (
	"log"

	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// func CreateAdminSeeder(db *gorm.DB) {
// 	var adminCount int64
// 	cfg := config.AppConfig

// 	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)

// 	if adminCount == 0 {
// 		log.Println("No admin user found. Seeding a new admin...")

// 		adminEmail := cfg.AdminEmail
// 		adminPassword := cfg.AdminPassword

// 		if adminEmail == "" || adminPassword == "" {
// 			log.Fatal("FATAL: ADMIN_EMAIL and ADMIN_PASSWORD must be set in .env to seed the first admin.")
// 		}

// 		newAdmin := models.User{
// 			Name:     "Default Admin",
// 			Email:    adminEmail,
// 			Password: adminPassword,
// 			Role:     models.RoleAdmin,
// 		}

// 		if err := db.Create(&newAdmin).Error; err != nil {
// 			log.Fatalf("FATAL: Failed to seed admin user: %v", err)
// 		}

// 		log.Printf("Successfully seeded admin user with email: %s\n", adminEmail)
// 	} else {
// 		log.Println("Admin user already exists. Seeder is skipping.")
// 	}
// }

func CreateAdminSeeder(db *gorm.DB, cfg config.Config) {
	var adminRole models.Role

	if err := db.Where(models.Role{Name: "admin"}).FirstOrCreate(&adminRole).Error; err != nil {
		log.Fatalf("FATAL: Failed to create admin role during seeding: %v", err)
	}
	log.Printf("Admin role ensured in database. Role ID is: %d", adminRole.ID)

	var userCount int64
	db.Model(&models.User{}).Where("email = ?", cfg.AdminEmail).Count(&userCount)
	if userCount > 0 {
		log.Println("Admin user already exists. Seeder is skipping.")
		return
	}

	log.Println("No admin user found. Seeding a new admin...")

	adminEmail := cfg.AdminEmail
	adminPassword := cfg.AdminPassword
	if adminEmail == "" || adminPassword == "" {
		log.Fatal("FATAL: Admin_Email and Admin_Password must be set in your .env file to seed the first admin.")
	}

	newAdmin := models.User{
		Name:     "Administrator",
		Email:    adminEmail,
		Password: adminPassword,
		Role:     models.RoleAdmin,
		// RoleID:   &adminRole.ID,
	}

	if err := db.Create(&newAdmin).Error; err != nil {
		log.Fatalf("FATAL: Failed to seed admin user: %v", err)
	}

	log.Printf("Successfully seeded admin user with email: %s and assigned RoleID: %d\n", adminEmail, adminRole.ID)
}
