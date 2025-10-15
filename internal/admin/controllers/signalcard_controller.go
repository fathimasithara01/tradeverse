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

	// ✅ Bind JSON body to request struct once
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Debug: Print the entire request data
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("📥 Received Signal Request:\n%s", string(reqJSON))

	// ✅ Normalize symbol name
	if !strings.HasSuffix(strings.ToUpper(req.Symbol), "USDT") {
		req.Symbol = strings.ToUpper(req.Symbol) + "USDT"
	} else {
		req.Symbol = strings.ToUpper(req.Symbol)
	}

	// ✅ Parse start and end dates
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		log.Printf("❌ Invalid startDate: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Must be RFC3339 (e.g., 2025-10-15T00:00:00Z)"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		log.Printf("❌ Invalid endDate: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Must be RFC3339 (e.g., 2025-10-16T00:00:00Z)"})
		return
	}

	// ✅ Map to Signal model
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
	}

	// ✅ Debug log specific fields
	log.Printf("✅ Parsed Signal Data: %+v", signal)

	// ✅ Save to DB through your service
	createdSignal, err := ctrl.liveSignalService.CreateSignal(c, &signal)
	if err != nil {
		log.Printf("❌ Error creating signal: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create signal"})
		return
	}

	// ✅ Success response
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
