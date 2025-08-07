package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type RevenueSplitRepository struct{}

func (r *RevenueSplitRepository) Save(split models.RevenueSplit) error {
	return db.DB.Create(&split).Error
}
