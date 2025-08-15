package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `gorm:"size:100;unique;not null" json:"name"`
	// CreatedByID uint   `json:"created_by_id"`
	// CreatedBy   User   `gorm:"foreignKey:CreatedByID"  json:"CreatedBy"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	// Users       []User       `gorm:"foreignKey:RoleID"`
}
