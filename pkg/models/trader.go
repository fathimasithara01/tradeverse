package models

import "gorm.io/gorm"

type TraderStatus string

const (
	StatusPending  TraderStatus = "pending"
	StatusApproved TraderStatus = "approved"
	StatusRejected TraderStatus = "rejected"
)

type TraderProfile struct {
	gorm.Model
	UserID      uint         `gorm:"unique;not null" json:"user_id"`
	CompanyName string       `gorm:"size:100" json:"company_name"`
	Bio         string       `json:"bio"`
	Status      TraderStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	TotalPnL    float64      `json:"total_pnl"`
	IsVerified  bool         `gorm:"default:false" json:"is_verified"`
}

// type TraderProfile struct {
// 	gorm.Model
// 	UserID          uint   `gorm:"unique;not null"`
// 	Bio             string `gorm:"type:text"`
// 	TradingStrategy string `gorm:"type:text"`
// 	ExperienceYears int    `gorm:"default:0"`
// 	IsVerified      bool   `gorm:"default:false"`
// 	// Add other trader-specific fields like performance metrics, links, etc.
// }
