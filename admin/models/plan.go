package models

import "gorm.io/gorm"

type Plan struct {
	gorm.Model
	TraderID    uint    `json:"trader_id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"` // in days
	IsActive    bool    `json:"is_active"`
}
