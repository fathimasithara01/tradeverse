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
	Name        string       `gorm:"size:100" json:"name"`
	CompanyName string       `gorm:"size:100" json:"company_name"`
	Bio         string       `json:"bio"`
	Status      TraderStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	TotalPnL    float64      `json:"total_pnl"`
	IsVerified  bool         `gorm:"default:false" json:"is_verified"`
}
