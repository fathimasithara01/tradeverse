package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedByID uint         `json:"created_by_id"`
	// CreatedBy   User         `gorm:"foreignKey:CreatedByID" json:"CreatedBy"`
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}
