package models

import (
	"time"

	"gorm.io/gorm"
)

type CopySession struct {
	gorm.Model
	FollowerID uint `gorm:"index"`
	Follower   User `gorm:"foreignKey:FollowerID"`
	MasterID   uint `gorm:"index"`
	Master     User `gorm:"foreignKey:MasterID"`

	RiskSetting   string  `gorm:"size:100"` // e.g., "Fixed Size: 0.01 lots"
	CurrentProfit float64 `gorm:"type:decimal(10,2)"`
	IsActive      bool    `gorm:"default:true;index"`
	StartedAt     time.Time
}
