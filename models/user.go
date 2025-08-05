package model

import (
"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleTrader Role = "trader"
	RoleUser Role = "user"
)

type User struct {
	gorm.Model
	Name string gorm:"not null"
	Email string gorm:"unique;not null"
	Password string gorm:"not null" // hashed
	Role Role gorm:"type:varchar(10);not null"
}