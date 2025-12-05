package seeder

import (
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

func CreateAdminSeeder(db *gorm.DB, cfg config.Config) {
	var adminUser models.User
	if db.Where("email = ?", cfg.Admin.Email).First(&adminUser).Error == gorm.ErrRecordNotFound {
		log.Println("Admin user not found. Seeding admin data...")

		adminRole := models.Role{Name: string(models.RoleAdmin), Description: "Super administrator with all permissions."}
		db.FirstOrCreate(&adminRole, models.Role{Name: string(models.RoleAdmin)})

		permissions := []models.Permission{
			{Name: "view_dashboard", Description: "View the admin dashboard"},
			{Name: "manage_users", Description: "Create, view, update, and delete users (customers and internal staff)"},
			{Name: "manage_roles", Description: "Create, view, update, and delete roles and assign permissions"},
			{Name: "manage_traders", Description: "Approve/reject trader applications and manage trader profiles"},
			{Name: "manage_signals", Description: "Create, update, and delete trading signals"},
			{Name: "view_activity_logs", Description: "View live copying sessions and trade error logs"},
			{Name: "manage_subscriptions", Description: "Manage subscription plans and user subscriptions"},
			{Name: "manage_wallet", Description: "Manage admin wallet, deposits, and withdrawals"},
			{Name: "view_transactions", Description: "View all platform transactions"},
			{Name: "delete_users", Description: "Permanently delete user accounts"}, 

			{Name: "view_admin_profile", Description: "View own admin profile details"},
			{Name: "edit_admin_profile", Description: "Edit own admin profile details, including password"},
			{Name: "view_admin_settings", Description: "View global admin settings"},
		}

		for _, perm := range permissions {
			db.FirstOrCreate(&perm, models.Permission{Name: perm.Name})
		}

		var createdPermissions []models.Permission
		db.Find(&createdPermissions)

		db.Model(&adminRole).Association("Permissions").Replace(createdPermissions)
		log.Printf("Admin role '%s' created and all permissions assigned.", adminRole.Name)

		admin := models.User{
			Name:      "Super Admin",
			Email:     cfg.Admin.Email,
			Password:  cfg.Admin.Password,
			Role:      models.RoleAdmin,
			RoleID:    &adminRole.ID,
			IsBlocked: false,
			IsVerified: true,
		}

		if err := admin.SetPassword(cfg.Admin.Password); err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
		log.Printf("Admin user '%s' created successfully.", admin.Email)

	} else if adminUser.ID != 0 {
		log.Printf("Admin user '%s' already exists. Skipping seeding.", cfg.Admin.Email)

		permissions := []models.Permission{
			{Name: "view_dashboard", Description: "View the admin dashboard"},
			{Name: "manage_users", Description: "Create, view, update, and delete users (customers and internal staff)"},
			{Name: "manage_roles", Description: "Create, view, update, and delete roles and assign permissions"},
			{Name: "manage_traders", Description: "Approve/reject trader applications and manage trader profiles"},
			{Name: "manage_signals", Description: "Create, update, and delete trading signals"},
			{Name: "view_activity_logs", Description: "View live copying sessions and trade error logs"},
			{Name: "manage_subscriptions", Description: "Manage subscription plans and user subscriptions"},
			{Name: "manage_wallet", Description: "Manage admin wallet, deposits, and withdrawals"},
			{Name: "view_transactions", Description: "View all platform transactions"},
			{Name: "delete_users", Description: "Permanently delete user accounts"},

			{Name: "view_admin_profile", Description: "View own admin profile details"},
			{Name: "edit_admin_profile", Description: "Edit own admin profile details, including password"},
			{Name: "view_admin_settings", Description: "View global admin settings"},
		}

		for _, perm := range permissions {
			db.FirstOrCreate(&perm, models.Permission{Name: perm.Name})
		}

		var adminRole models.Role
		db.Where("name = ?", models.RoleAdmin).First(&adminRole)
		if adminRole.ID != 0 {
			var existingPermissions []models.Permission
			db.Model(&adminRole).Association("Permissions").Find(&existingPermissions)

			existingPermNames := make(map[string]bool)
			for _, p := range existingPermissions {
				existingPermNames[p.Name] = true
			}

			var permissionsToAdd []models.Permission
			for _, p := range permissions {
				if !existingPermNames[p.Name] {
					var newPerm models.Permission
					db.Where("name = ?", p.Name).First(&newPerm)
					if newPerm.ID != 0 {
						permissionsToAdd = append(permissionsToAdd, newPerm)
					}
				}
			}

			if len(permissionsToAdd) > 0 {
				db.Model(&adminRole).Association("Permissions").Append(permissionsToAdd)
				log.Printf("Added %d new permissions to admin role.", len(permissionsToAdd))
			} else {
				log.Println("Admin role already has all necessary permissions.")
			}
		}
	}
}
