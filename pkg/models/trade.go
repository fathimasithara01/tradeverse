package models

import (
	"time"

	"gorm.io/gorm"
)

type Trade struct {
	gorm.Model
	MasterUserID  uint    `gorm:"index"`
	Symbol        string  `gorm:"size:20;not null"`
	EntryPrice    float64 `gorm:"not null"`
	StopLossPrice float64
	TargetPrice   float64
	Status        string    `gorm:"size:20;default:'open';index"`
	OpenedAt      time.Time `gorm:"not null"`
}
