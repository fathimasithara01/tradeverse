package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name     string `gorm:"size:100;unique;not null" json:"name"`
	Category string `gorm:"size:100;not null" json:"category"`
}
