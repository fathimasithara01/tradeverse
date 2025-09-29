package service

// import (
// 	"errors"

// 	"github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// type TraderWalletService struct {
// 	db         *gorm.DB
// 	walletRepo walletrepo.TraderWalletRepository
// }

// func NewTraderWalletService(db *gorm.DB, walletRepo walletrepo.TraderWalletRepository) *TraderWalletService {
// 	return &TraderWalletService{db: db, walletRepo: walletRepo}
// }

// func (s *TraderWalletService) SubscribeCustomer(customerID, traderID uint, price float64, currency string) error {
// 	return s.db.Transaction(func(tx *gorm.DB) error {
// 		customerWallet, err := s.walletRepo.GetByUserID(customerID)
// 		if err != nil {
// 			return errors.New("customer wallet not found")
// 		}
// 		traderWallet, err := s.walletRepo.GetByUserID(traderID)
// 		if err != nil {
// 			return errors.New("trader wallet not found")
// 		}
// 		adminWallet, err := s.walletRepo.GetByUserID(1)
// 		if err != nil {
// 			return errors.New("admin wallet not found")
// 		}

// 		if customerWallet.Balance < price {
// 			return errors.New("insufficient balance")
// 		}

// 		adminShare := price * 0.20
// 		traderShare := price * 0.80

// 		if err := s.walletRepo.UpdateBalance(customerWallet.ID, -price); err != nil {
// 			return err
// 		}
// 		_ = s.walletRepo.AddTransaction(&models.WalletTransaction{
// 			WalletID:        customerWallet.ID,
// 			UserID:          customerID,
// 			TransactionType: models.TxTypeSubscription,
// 			Amount:          -price,
// 			Currency:        currency,
// 			Status:          models.TxStatusSuccess,
// 			Description:     "subscription payment",
// 		})

// 		if err := s.walletRepo.UpdateBalance(adminWallet.ID, adminShare); err != nil {
// 			return err
// 		}
// 		_ = s.walletRepo.AddTransaction(&models.WalletTransaction{
// 			WalletID:        adminWallet.ID,
// 			UserID:          1,
// 			TransactionType: models.TxTypeSubscription,
// 			Amount:          adminShare,
// 			Currency:        currency,
// 			Status:          models.TxStatusSuccess,
// 			Description:     "subscription share",
// 		})

// 		if err := s.walletRepo.UpdateBalance(traderWallet.ID, traderShare); err != nil {
// 			return err
// 		}
// 		_ = s.walletRepo.AddTransaction(&models.WalletTransaction{
// 			WalletID:        traderWallet.ID,
// 			UserID:          traderID,
// 			TransactionType: models.TxTypeSubscription,
// 			Amount:          traderShare,
// 			Currency:        currency,
// 			Status:          models.TxStatusSuccess,
// 			Description:     "subscription earnings",
// 		})

// 		return nil
// 	})
// }

// func (s *TraderWalletService) GetBalance(userID uint) (*models.WalletSummaryResponse, error) {
// 	wallet, err := s.walletRepo.GetByUserID(userID)
// 	if err != nil {
// 		return nil, errors.New("wallet not found")
// 	}

// 	return &models.WalletSummaryResponse{
// 		UserID:      userID,
// 		WalletID:    wallet.ID,
// 		Balance:     wallet.Balance,
// 		Currency:    wallet.Currency,
// 		LastUpdated: wallet.UpdatedAt,
// 	}, nil
// }
