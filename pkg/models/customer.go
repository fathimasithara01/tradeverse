package models

import "gorm.io/gorm"

type CustomerProfile struct {
	gorm.Model
	Name            string `json:"name"`
	UserID          uint   `gorm:"unique;not null"`
	ShippingAddress string `gorm:"size:255"`
	Phone           string `gorm:"size:20"`
}
