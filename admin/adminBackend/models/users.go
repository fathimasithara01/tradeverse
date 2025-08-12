package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole defines a clear, controlled set of roles for users.
type UserRole string

// Constants provide type-safety and prevent typos for role names.
const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
	RoleTrader   UserRole = "trader"
)

// User is the central model for any entity that can log in.
// It contains the core information for authentication and authorization.
type User struct {
	gorm.Model          // Standard GORM fields: ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string   `gorm:"size:100;not null" json:"name"`
	Email      string   `gorm:"size:100;unique;not null" json:"email"`
	Password   string   `gorm:"size:255;not null" json:"-"` // json:"-" PREVENTS the password from ever being sent in an API response
	Role       UserRole `gorm:"type:varchar(20);not null" json:"role"`
	IsBlocked  bool     `gorm:"default:false" json:"is_blocked"`

	CustomerProfile CustomerProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // A user can have one customer profile
	TraderProfile   TraderProfile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // A user can have one trader profile
}

// BeforeCreate is a GORM hook to automatically hash passwords for new users.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}
