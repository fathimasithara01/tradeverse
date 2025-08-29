package models

import "gorm.io/gorm"

// Product represents an item or service in your system.
type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
