package controllers

// import (
// 	"net/http"

// 	"github.com/fathimasithara01/tradeverse/admin/service"
// 	"github.com/gin-gonic/gin"
// )

// func ExportUsersCSV(c *gin.Context) {
// 	fileData, err := service.ExportService{}.GenerateUsersCSV()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export CSV"})
// 		return
// 	}

// 	c.Header("Content-Disposition", "attachment; filename=users.csv")
// 	c.Data(http.StatusOK, "text/csv", fileData)
// }

// func ExportSubscriptionsPDF(c *gin.Context) {
// 	pdfBytes, err := service.ExportService{}.GenerateSubscriptionsPDF()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
// 		return
// 	}
// 	c.Header("Content-Disposition", "attachment; filename=subscriptions.pdf")
// 	c.Data(http.StatusOK, "application/pdf", pdfBytes)
// }
