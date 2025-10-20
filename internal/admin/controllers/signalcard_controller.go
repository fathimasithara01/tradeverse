package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignalController struct {
	liveSignalService service.ILiveSignalService
}

func NewSignalController(liveSignalService service.ILiveSignalService) *SignalController {
	return &SignalController{
		liveSignalService: liveSignalService,
	}
}

func (ctrl *SignalController) ShowLiveSignalsPage(c *gin.Context) {
	signals, err := ctrl.liveSignalService.GetAllSignals(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signal_cards.html", gin.H{
			"Title":        "Live Trading Signals",
			"ActiveTab":    "activity",
			"ActiveSubTab": "signal_cards",
			"Error":        "Failed to load signals",
		})
		return
	}

	c.HTML(http.StatusOK, "signal_cards.html", gin.H{
		"Title":        "Live Trading Signals",
		"ActiveTab":    "activity",
		"ActiveSubTab": "signal_cards",
		"Signals":      signals,
	})
}

func (ctrl *SignalController) GetLiveSignals(c *gin.Context) {
	signals, err := ctrl.liveSignalService.GetAllSignals(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve signals"})
		return
	}
	c.JSON(http.StatusOK, signals)
}

func GetSignalCardsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signal_cards.html", gin.H{
		"Title":        "Signal Cards",
		"ActiveTab":    "activity",
		"ActiveSubTab": "signal_cards",
	})
}

func (ctrl *SignalController) ShowCreateSignalCardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_signal_card.html", gin.H{
		"Title":        "Create Signal Card",
		"ActiveTab":    "activity",
		"ActiveSubTab": "signal_cards",
	})
}

// ... (imports)

// You'll need to modify this function in your controllers/signal_controller.go
func (ctrl *SignalController) CreateSignal(c *gin.Context) {
	var req struct {
		TraderName    string  `json:"traderName"`
		Symbol        string  `json:"symbol"`
		StopLoss      float64 `json:"stopLoss"`
		EntryPrice    float64 `json:"entryPrice"`
		TargetPrice   float64 `json:"targetPrice"`
		CurrentPrice  float64 `json:"currentPrice"`
		Risk          string  `json:"risk"`
		Strategy      string  `json:"strategy"`
		Status        string  `json:"status"`
		StartDate     string  `json:"startDate"`
		EndDate       string  `json:"endDate"`
		TotalDuration string  `json:"totalDuration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("Received Signal Request:\n%s", string(reqJSON))

	if !strings.HasSuffix(strings.ToUpper(req.Symbol), "USDT") {
		req.Symbol = strings.ToUpper(req.Symbol) + "USDT"
	} else {
		req.Symbol = strings.ToUpper(req.Symbol)
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		log.Printf(" Invalid startDate: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Must be RFC3339 (e.g., 2025-10-15T00:00:00Z)"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		log.Printf(" Invalid endDate: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Must be RFC3339 (e.g., 2025-10-16T00:00:00Z)"})
		return
	}

	// --- NEW LOGIC HERE ---
	// You'll need to get the current user's role and ID from your authentication system.
	// This is a placeholder; replace with your actual authentication logic.
	var createdByRole string
	var creatorID uint

	// Example: Assuming user data is stored in c.MustGet("user") after middleware
	// This part heavily depends on how your auth middleware sets user data.
	/*
		if user, exists := c.Get("user"); exists {
			if adminUser, ok := user.(*models.AdminUser); ok { // Assuming you have an AdminUser model
				createdByRole = "Admin"
				creatorID = adminUser.ID
				// For admin creating a signal, the traderName might be from the form or selected admin
				if req.TraderName == "" {
					req.TraderName = adminUser.Username // or a default admin trader name
				}
			} else if traderUser, ok := user.(*models.TraderUser); ok { // Assuming you have a TraderUser model
				createdByRole = "Trader"
				creatorID = traderUser.ID
				req.TraderName = traderUser.Username // Trader creating a signal, their own name is the trader name
			}
		} else {
			// Handle unauthenticated request or default
			createdByRole = "Unknown"
			creatorID = 0
		}
	*/

	// For demonstration, let's hardcode for the admin panel for now.
	// In a real application, this would be dynamic.
	// Since this is specifically for the admin/create_signal_card.html,
	// we can assume an admin is creating it.
	createdByRole = "Admin"
	// You might get the admin's ID from the session or JWT
	creatorID = 1 // Placeholder for Admin ID

	signal := models.Signal{
		TraderName:     req.TraderName,
		Symbol:         req.Symbol,
		StopLoss:       req.StopLoss,
		EntryPrice:     req.EntryPrice,
		TargetPrice:    req.TargetPrice,
		CurrentPrice:   req.CurrentPrice,
		Risk:           req.Risk,
		Strategy:       req.Strategy,
		Status:         req.Status,
		TotalDuration:  req.TotalDuration,
		TradeStartDate: startDate,
		TradeEndDate:   endDate,
		PublishedAt:    time.Now(),
		CreatedBy:      createdByRole, // Set who created it
		CreatorID:      creatorID,     // Set the ID of the creator
	}
	
	log.Printf(" Parsed Signal Data: %+v", signal)

	createdSignal, err := ctrl.liveSignalService.CreateSignal(c, &signal)
	if err != nil {
		log.Printf("Error creating signal: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create signal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signal created successfully",
		"data":    createdSignal,
	})
}

func GetMarketDataAPI(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found in context"})
		log.Println("ERROR: Database connection not found in Gin context for GetMarketDataAPI")
		return
	}
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database context object is not a GORM DB instance"})
		log.Println("ERROR: Database context object is not a GORM DB instance for GetMarketDataAPI")
		return
	}

	var marketData []models.MarketData
	log.Println("Attempting to retrieve market data from DB for /admin/api/market-data...")
	if err := gormDB.Order("current_price DESC").Find(&marketData).Error; err != nil {
		log.Printf("ERROR: Failed to retrieve market data from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve market data"})
		return
	}

	log.Printf("Successfully retrieved %d market data entries for /admin/api/market-data.", len(marketData))
	var apiResponse []models.MarketDataAPIResponse
	for _, md := range marketData {
		apiResponse = append(apiResponse, models.MarketDataAPIResponse{
			Symbol:         md.Symbol,
			Name:           md.Name,
			CurrentPrice:   md.CurrentPrice,
			PriceChange24H: md.PriceChange24H,
			LogoURL:        md.LogoURL,
		})
	}

	c.JSON(http.StatusOK, apiResponse)
}
