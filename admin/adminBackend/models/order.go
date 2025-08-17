package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	// CORRECTED: The foreign key must be to the 'users' table.
	CustomerID uint `json:"customer_id"`
	Customer   User `gorm:"foreignKey:CustomerID"` // An Order "Belongs To" a User

	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`

	ShippingAddressID uint    `json:"shipping_address_id"`
	ShippingAddress   Address `gorm:"foreignKey:ShippingAddressID"`

	OrderDate     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"order_date"`
	TotalAmount   float64   `gorm:"type:decimal(10,2)" json:"total_amount"`
	Status        string    `gorm:"default:'Pending'" json:"status"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `gorm:"default:'Pending'" json:"payment_status"`
	PaymentID     string    `json:"payment_id,omitempty"`
}

type OrderItem struct {
	gorm.Model
	OrderID         uint    `json:"order_id"`
	ProductID       uint    `json:"product_id"`
	Product         Product `gorm:"foreignKey:ProductID"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `gorm:"type:decimal(10,2)" json:"price_at_purchase"`
}

type Address struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}
