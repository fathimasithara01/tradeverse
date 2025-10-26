package models

import "gorm.io/gorm"

type WebConfiguration struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex;not null" json:"key"`
	Value string `gorm:"not null" json:"value"`
	Type  string `gorm:"not null;default:'string'" json:"type"` // e.g., 'string', 'number', 'boolean', 'json'
}
