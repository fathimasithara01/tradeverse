package models

import (
	"time"

	"gorm.io/gorm"
)

type LogStatus string

const (
	LogStatusSuccess LogStatus = "success"
	LogStatusFailed  LogStatus = "failed"
)

type TradeLog struct {
	gorm.Model
	CopySessionID   uint
	MasterTradeID   string
	FollowerTradeID string
	Status          LogStatus `gorm:"type:varchar(20);index"`
	ErrorMessage    string    
	ExecutionTime   int64  
	Timestamp       time.Time
}
