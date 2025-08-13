package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
	RoleTrader   UserRole = "trader"
)

type User struct {
	gorm.Model
	Name      string   `gorm:"size:100;not null" json:"name"`
	Email     string   `gorm:"size:100;unique;not null" json:"email"`
	Password  string   `gorm:"size:255;not null" json:"-"`
	Role      UserRole `gorm:"type:varchar(20);not null" json:"role"`
	IsBlocked bool     `gorm:"default:false" json:"is_blocked"`

	RolesCreated []Role `gorm:"foreignKey:CreatedByID"`

	CustomerProfile CustomerProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	TraderProfile   TraderProfile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}
