package service

// import (
// 	"bytes"
// 	"encoding/csv"
// 	"fmt"

// 	"github.com/fathimasithara01/tradeverse/admin/db"
// 	"github.com/fathimasithara01/tradeverse/admin/models"
// 	"github.com/jung-kurt/gofpdf"
// )

// type ExportService struct{}

// //  Generate CSV of all users
// func (s *ExportService) GenerateUsersCSV() ([]byte, error) {
// 	var users []models.User
// 	if err := db.DB.Find(&users).Error; err != nil {
// 		return nil, err
// 	}

// 	var buffer bytes.Buffer
// 	writer := csv.NewWriter(&buffer)

// 	// Headers
// 	writer.Write([]string{"ID", "Name", "Email", "Role", "Status"})

// 	for _, u := range users {
// 		writer.Write([]string{
// 			fmt.Sprint(u.ID),
// 			u.Name,
// 			u.Email,
// 			u.Role,
// 			u.Status,
// 		})
// 	}
// 	writer.Flush()
// 	return buffer.Bytes(), nil
// }

// // ✅ Generate PDF of all subscriptions
// func (s *ExportService) GenerateSubscriptionsPDF() ([]byte, error) {
// 	var subscriptions []models.Subscription
// 	if err := db.DB.Preload("User").Preload("Plan").Find(&subscriptions).Error; err != nil {
// 		return nil, err
// 	}

// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()
// 	pdf.SetFont("Arial", "B", 14)
// 	pdf.Cell(40, 10, "Subscription Report")
// 	pdf.Ln(12)

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.Cell(20, 10, "ID")
// 	pdf.Cell(40, 10, "User")
// 	pdf.Cell(40, 10, "Plan")
// 	pdf.Cell(30, 10, "Status")
// 	pdf.Ln(10)

// 	for _, s := range subscriptions {
// 		pdf.Cell(20, 10, fmt.Sprint(s.ID))
// 		pdf.Cell(40, 10, s.User.Name)
// 		pdf.Cell(40, 10, s.Plan.Name)
// 		pdf.Cell(30, 10, s.Status)
// 		pdf.Ln(8)
// 	}

// 	var buf bytes.Buffer
// 	err := pdf.Output(&buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }
