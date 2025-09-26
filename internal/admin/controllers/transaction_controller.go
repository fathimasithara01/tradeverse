package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service service.TransactionService
}

func NewTransactionController(service service.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

func (c *TransactionController) GetTransactionsAPI(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	year, _ := strconv.Atoi(ctx.Query("year"))
	month, _ := strconv.Atoi(ctx.Query("month"))
	day, _ := strconv.Atoi(ctx.Query("day"))

	transactions, total, err := c.service.GetTransactions(page, limit, search, year, month, day)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

func (c *TransactionController) GetTransactionsPage(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	selectedYear, _ := strconv.Atoi(ctx.DefaultQuery("year", "0"))  
	selectedMonth, _ := strconv.Atoi(ctx.DefaultQuery("month", "0")) 
	selectedDay, _ := strconv.Atoi(ctx.DefaultQuery("day", "0"))  

	transactions, total, err := c.service.GetTransactions(page, limit, search, selectedYear, selectedMonth, selectedDay)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Failed to fetch transactions"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Generate page numbers for iteration
	pages := make([]int, totalPages)
	for i := 0; i < totalPages; i++ {
		pages[i] = i + 1
	}

	availableYears, err := c.service.GetAvailableYears()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available years"})
		return
	}

	monthNames := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	daysInMonth := 31
	days := make([]int, daysInMonth)
	for i := 0; i < daysInMonth; i++ {
		days[i] = i + 1
	}

	ctx.HTML(http.StatusOK, "transactions.html", gin.H{
		"transactions":  transactions,
		"total":         total,
		"page":          page,
		"limit":         limit,
		"search":        search,
		"totalPages":    totalPages,
		"pages":         pages,
		"years":         availableYears,
		"months":        monthNames,
		"days":          days,
		"selectedYear":  selectedYear,
		"selectedMonth": selectedMonth,
		"selectedDay":   selectedDay,
		"ActiveTab":     "financials",  
		"ActiveSubTab":  "transactions", 
	})
}
