package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
	Type    string `json:"type"` // info, alert, warning
}
