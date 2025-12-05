package wallet

import "errors"

type Wallet struct {
	balance int
}

func NewWallet(initial int) *Wallet {
	return &Wallet{balance: initial}
}

func (w *Wallet) Deposit(amount int) error {
	if amount <= 0 {
		return errors.New("invalid deposit amount")
	}
	w.balance += amount
	return nil
}

func (w *Wallet) Withdraw(amount int) error {
	if amount <= 0 {
		return errors.New("invalid withdraw amount")
	}
	if amount > w.balance {
		return errors.New("insufficient balance")
	}
	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() int {
	return w.balance
}
