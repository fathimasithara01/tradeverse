package tests

import (
	"testing"

	"github.com/fathimasithara01/tradeverse/pkg/wallet"
)

func TestWalletDeposit(t *testing.T) {
	w := wallet.NewWallet(0)

	err := w.Deposit(100)
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}

	if w.Balance() != 100 {
		t.Errorf("expected balance 100, got %d", w.Balance())
	}
}

func TestWalletWithdraw(t *testing.T) {
	w := wallet.NewWallet(200)

	err := w.Withdraw(50)
	if err != nil {
		t.Fatalf("withdraw failed: %v", err)
	}

	if w.Balance() != 150 {
		t.Errorf("expected balance 150, got %d", w.Balance())
	}
}
