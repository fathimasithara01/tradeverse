package seeder

import (
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateAdminSeeder(db *gorm.DB, cfg config.Config) {
	// Ensure Admin role
	var adminRole models.Role
	if err := db.Where(models.Role{Name: string(models.RoleAdmin)}).FirstOrCreate(&adminRole).Error; err != nil {
		log.Fatalf("FATAL: Failed to create or find admin role during seeding: %v", err)
	}
	log.Printf("Admin role ensured in database. Role ID is: %d", adminRole.ID)

	// Ensure Customer role
	var customerRole models.Role
	if err := db.Where(models.Role{Name: string(models.RoleCustomer)}).FirstOrCreate(&customerRole).Error; err != nil {
		log.Fatalf("FATAL: Failed to create or find customer role during seeding: %v", err)
	}
	log.Printf("Customer role ensured in database. Role ID is: %d", customerRole.ID)

	// Ensure Trader role
	var traderRole models.Role
	if err := db.Where(models.Role{Name: string(models.RoleTrader)}).FirstOrCreate(&traderRole).Error; err != nil {
		log.Fatalf("FATAL: Failed to create or find trader role during seeding: %v", err)
	}
	log.Printf("Trader role ensured in database. Role ID is: %d", traderRole.ID)

	// Seed Admin user
	var userCount int64
	db.Model(&models.User{}).Where("email = ?", cfg.Admin.Email).Count(&userCount)
	if userCount > 0 {
		log.Println("Admin user already exists. Seeder is skipping.")
		return
	}

	log.Println("No admin user found. Seeding a new admin...")

	adminEmail := cfg.Admin.Email
	adminPassword := cfg.Admin.Password
	if adminEmail == "" || adminPassword == "" {
		log.Fatal("FATAL: Admin_Email and Admin_Password must be set in your .env file to seed the first admin.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("FATAL: Failed to hash admin password: %v", err)
	}

	newAdmin := models.User{
		Name:     "Administrator",
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     models.RoleAdmin,
		RoleID:   &adminRole.ID,
	}

	if err := db.Create(&newAdmin).Error; err != nil {
		log.Fatalf("FATAL: Failed to seed admin user: %v", err)
	}

	log.Printf("Successfully seeded admin user with email: %s and assigned RoleID: %d\n", adminEmail, adminRole.ID)
}
