package models

import "gorm.io/gorm"

type Follower struct {
	gorm.Model
	TraderID uint
	UserID   uint
}
