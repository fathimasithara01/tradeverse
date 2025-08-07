package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/razorpay/razorpay-go"
)

type PaymentService struct {
	Repo           repository.PaymentRepository
	RazorpayKey    string
	RazorpaySecret string
}

func (s *PaymentService) GetAll() ([]models.Payment, error) {
	return s.Repo.GetAllPayments()
}

func (s *PaymentService) CreateRazorpayOrder(userID uint, amount int64) (string, error) {
	client := razorpay.NewClient(s.RazorpayKey, s.RazorpaySecret)

	data := map[string]interface{}{
		"amount":   amount * 100, // in paisa
		"currency": "INR",
		"receipt":  "rcpt_" + string(userID),
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}

	orderID := body["id"].(string)
	payment := models.Payment{
		UserID:  userID,
		OrderID: orderID,
		Amount:  float64(amount),
		Method:  "razorpay",
		Status:  "pending",
	}
	_ = s.Repo.Save(payment)

	return orderID, nil
}

func (s *PaymentService) VerifyAndSave(userID, traderID uint, paymentID, orderID, signature string, amount int64) error {
	// 1. Save payment
	payment := models.Payment{
		UserID:    userID,
		TraderID:  traderID,
		OrderID:   orderID,
		PaymentID: paymentID,
		Amount:    float64(amount),
		Status:    "success",
		Method:    "razorpay",
	}
	err := s.Repo.Save(payment)
	if err != nil {
		return err
	}

	// 2. Revenue Split Calculation
	adminShare := float64(amount) * 0.2  // 20%
	traderShare := float64(amount) * 0.8 // 80%

	split := models.RevenueSplit{
		PaymentID:   payment.ID,
		UserID:      userID,
		TraderID:    traderID,
		AdminShare:  adminShare,
		TraderShare: traderShare,
		TotalAmount: float64(amount),
	}

	return (&repository.RevenueSplitRepository{}).Save(split)
}

func (s *PaymentService) GetAllPayments() ([]models.Payment, error) {
	return s.Repo.GetAll()
}
