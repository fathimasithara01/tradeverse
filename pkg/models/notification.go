package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeTradeUpdate    NotificationType = "TRADE_UPDATE"
	NotificationTypeSubscription   NotificationType = "SUBSCRIPTION"
	NotificationTypeWallet         NotificationType = "WALLET"
	NotificationTypePlatform       NotificationType = "PLATFORM_ANNOUNCEMENT"
	NotificationTypeCopierActivity NotificationType = "COPIER_ACTIVITY"
)

type Notification struct {
	gorm.Model
	UserID    uint             `gorm:"not null;index" json:"user_id"` // Recipient of the notification
	User      User             `gorm:"foreignKey:UserID" json:"user"`
	Type      NotificationType `gorm:"size:50;not null" json:"type"`
	Title     string           `gorm:"size:255;not null" json:"title"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Read      bool             `gorm:"default:false" json:"read"`
	Link      string           `gorm:"size:255" json:"link,omitempty"` // Optional link to related resource
	Timestamp time.Time        `gorm:"not null" json:"timestamp"`
}