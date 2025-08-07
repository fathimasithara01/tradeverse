package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type PaymentRepository struct{}

func (r *PaymentRepository) GetAllPayments() ([]models.Payment, error) {
	var payments []models.Payment
	err := db.DB.Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) Save(payment models.Payment) error {
	return db.DB.Create(&payment).Error
}

func (r *PaymentRepository) GetAll() ([]models.Payment, error) {
	var payments []models.Payment
	err := db.DB.Find(&payments).Error
	return payments, err
}
