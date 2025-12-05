package models

import (
	"gorm.io/gorm"
)

type WebConfiguration struct {
	gorm.Model
	PrimaryCountry   string `gorm:"type:varchar(100);not null;default:'United Arab Emirates'" json:"primary_country"`
	PrimaryCurrency  string `gorm:"type:varchar(100);not null;default:'United Arab Emirates Dirham (AED)'" json:"primary_currency"` // Increased size
	PrimaryTimezone  string `gorm:"type:varchar(100);not null;default:'Asia/Dubai'" json:"primary_timezone"`                        // e.g., Asia/Dubai, America/New_York
	FilesystemConfig string `gorm:"type:text" json:"filesystem_config"`                                                             // Placeholder for filesystem settings
}

func EnsureDefaultWebConfiguration(db *gorm.DB) error {
	var count int64
	db.Model(&WebConfiguration{}).Count(&count)
	if count == 0 {
		defaultConfig := WebConfiguration{
			PrimaryCountry:   "United Arab Emirates",
			PrimaryCurrency:  "United Arab Emirates Dirham (AED)", 
			PrimaryTimezone:  "Asia/Dubai",
			FilesystemConfig: "System", 
		}
		if err := db.Create(&defaultConfig).Error; err != nil {
			return err
		}
	}
	return nil
}
