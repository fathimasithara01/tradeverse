// models/web_configuration.go
package models

import (
	"gorm.io/gorm"
)

// WebConfiguration stores the general web configuration settings.
type WebConfiguration struct {
	gorm.Model
	PrimaryCountry   string `gorm:"type:varchar(100);not null;default:'United Arab Emirates'" json:"primary_country"`
	PrimaryCurrency  string `gorm:"type:varchar(10);not null;default:'AED'" json:"primary_currency"`         // e.g., AED, USD, EUR
	PrimaryTimezone  string `gorm:"type:varchar(100);not null;default:'Asia/Dubai'" json:"primary_timezone"` // e.g., Asia/Dubai, America/New_York
	FilesystemConfig string `gorm:"type:text" json:"filesystem_config"`                                      // Placeholder for filesystem settings
	// Add other web configuration settings here as needed
}

// EnsureDefaultWebConfiguration checks if a default web configuration exists and creates one if not.
func EnsureDefaultWebConfiguration(db *gorm.DB) error {
	var count int64
	db.Model(&WebConfiguration{}).Count(&count)
	if count == 0 {
		defaultConfig := WebConfiguration{
			PrimaryCountry:   "United Arab Emirates",
			PrimaryCurrency:  "United Arab Emirates Dirham (AED)", // Changed for better display
			PrimaryTimezone:  "Asia/Dubai",
			FilesystemConfig: "System", // Default value
		}
		if err := db.Create(&defaultConfig).Error; err != nil {
			return err
		}
	}
	return nil
}
