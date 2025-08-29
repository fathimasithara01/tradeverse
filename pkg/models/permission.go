package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`

	Category string `gorm:"size:100;not null;index" json:"category"`
}
