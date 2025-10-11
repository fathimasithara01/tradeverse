package service

import (
	"context"
	"errors"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type WalletService interface {
	GetBalance(ctx context.Context, userID uint) (*models.Wallet, error)
	Deposit(ctx context.Context, userID uint, amount float64) (*models.WalletTransaction, error)
	Withdraw(ctx context.Context, userID uint, amount float64) (*models.WalletTransaction, error)
	GetTransactionHistory(ctx context.Context, userID uint) ([]models.WalletTransaction, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetBalance(ctx context.Context, userID uint) (*models.Wallet, error) {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}
	return wallet, nil
}

func (s *walletService) Deposit(ctx context.Context, userID uint, amount float64) (*models.WalletTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	wallet.Balance += amount
	if err := s.repo.UpdateWallet(ctx, wallet); err != nil {
		return nil, err
	}

	tx := &models.WalletTransaction{
		WalletID: wallet.ID,
		UserID:   userID,
		// Type:     models.Deposit,
		Amount: amount,
	}
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *walletService) Withdraw(ctx context.Context, userID uint, amount float64) (*models.WalletTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}
	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	wallet.Balance -= amount
	if err := s.repo.UpdateWallet(ctx, wallet); err != nil {
		return nil, err
	}

	tx := &models.WalletTransaction{
		WalletID: wallet.ID,
		UserID:   userID,
		// Type:     models.Withdrawal,
		Amount: amount,
	}
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *walletService) GetTransactionHistory(ctx context.Context, userID uint) ([]models.WalletTransaction, error) {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}
	return s.repo.GetTransactionsByWalletID(ctx, wallet.ID)
}
