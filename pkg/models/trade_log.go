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
	ErrorMessage    string    // Populated only if Status is "failed"
	ExecutionTime   int64     // Latency in milliseconds
	Timestamp       time.Time
}
