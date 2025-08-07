package models

import "gorm.io/gorm"

type Announcement struct {
	gorm.Model
	Title   string `json:"title"`
	Message string `json:"message"`
	Target  string `json:"target"` // "all", "trader", "user"
}
