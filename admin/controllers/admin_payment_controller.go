package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var paymentService = service.PaymentService{
	Repo:           repository.PaymentRepository{},
	RazorpayKey:    "rzp_test_xxxxxxxx", // your key here
	RazorpaySecret: "xxxxxxxxxxxxxxxxx", // your secret here
}

func GetAllPayments(c *gin.Context) {
	payments, err := paymentService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment records"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func CreateCheckout(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	var body struct {
		Amount int64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	orderID, err := paymentService.CreateRazorpayOrder(uint(userID), body.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_id": orderID})
}

func VerifyPayment(c *gin.Context) {
	var body struct {
		UserID    uint   `json:"user_id"`
		TraderID  uint   `json:"trader_id"`
		OrderID   string `json:"order_id"`
		PaymentID string `json:"payment_id"`
		Signature string `json:"signature"`
		Amount    int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := paymentService.VerifyAndSave(body.UserID, body.TraderID, body.PaymentID, body.OrderID, body.Signature, body.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment verified & revenue split recorded"})
}

func AdminPaymentHistory(c *gin.Context) {
	payments, err := paymentService.GetAllPayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fetch failed"})
		return
	}
	c.JSON(http.StatusOK, payments)
}
