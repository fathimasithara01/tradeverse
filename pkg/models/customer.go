package models

import "gorm.io/gorm"

type CustomerProfile struct {
	gorm.Model
	UserID          uint   `gorm:"unique;not null"`
	ShippingAddress string `gorm:"size:255"`
	PhoneNumber     string `gorm:"size:20"`
}
