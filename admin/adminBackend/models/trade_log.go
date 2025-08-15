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

// TradeLog records every single attempt to replicate a trade.
type TradeLog struct {
	gorm.Model
	CopySessionID   uint
	MasterTradeID   string // The unique ID of the original trade from the master's broker
	FollowerTradeID string // The ID of the replicated trade on the follower's broker

	Status        LogStatus `gorm:"type:varchar(20);index"`
	ErrorMessage  string    // Populated only if Status is "failed"
	ExecutionTime int64     // Latency in milliseconds
	Timestamp     time.Time
}
