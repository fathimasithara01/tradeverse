package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type WalletRepository struct{}

func (r *WalletRepository) GetWalletByUserID(userID uint) (models.Wallet, []models.WalletTransaction, error) {
	var wallet models.Wallet
	var txs []models.WalletTransaction

	err := db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID}).Error
	_ = db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&txs).Error // ignore tx error

	// 	err1 := db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID}).Error
	// err2 := db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&txs).Error

	// if err1 != nil || err2 != nil {
	// 	return wallet, txs, fmt.Errorf("wallet error: %v, tx error: %v", err1, err2)
	// }
	return wallet, txs, err
}

func (r *WalletRepository) Credit(userID uint, amount float64, desc string) error {
	var wallet models.Wallet
	db.DB.FirstOrCreate(&wallet, models.Wallet{UserID: userID})

	wallet.Balance += amount
	if err := db.DB.Save(&wallet).Error; err != nil {
		return err
	}

	tx := models.WalletTransaction{
		UserID:      userID,
		Type:        "credit",
		Amount:      amount,
		Description: desc,
	}
	return db.DB.Create(&tx).Error
}

func (r *WalletRepository) Debit(userID uint, amount float64, desc string) error {
	var wallet models.Wallet
	db.DB.FirstOrCreate(&wallet, models.Wallet{UserID: userID})

	if wallet.Balance < amount {
		return nil // or return an error indicating insufficient funds
	}
	wallet.Balance -= amount
	if err := db.DB.Save(&wallet).Error; err != nil {
		return err
	}

	tx := models.WalletTransaction{
		UserID:      userID,
		Type:        "debit",
		Amount:      amount,
		Description: desc,
	}
	return db.DB.Create(&tx).Error
}
