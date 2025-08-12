package models

import "gorm.io/gorm"

type TraderProfile struct {
	gorm.Model
	UserID      uint   `gorm:"unique;not null"`
	CompanyName string `gorm:"size:100"`
	Bio         string
	TotalPnL    float64
	IsVerified  bool `gorm:"default:false"`
}
