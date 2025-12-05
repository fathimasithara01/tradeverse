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
	c.HTML(http.StatusOK, "signal_cards.html", gin.H{
		"Title":        "Live Trading Signals",
		"ActiveTab":    "activity",
		"ActiveSubTab": "signal_cards",
	})
}

func (ctrl *SignalController) GetLiveSignals(c *gin.Context) {
	signals, err := ctrl.liveSignalService.GetAllSignals(c)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve signals in GetLiveSignals controller: %v", err)
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

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: Failed to bind JSON for CreateSignal: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("Received Signal Request:\n%s", string(reqJSON))

	normalizedSymbol := strings.ToUpper(req.Symbol)
	if !strings.HasSuffix(normalizedSymbol, "USDT") {
		normalizedSymbol += "USDT"
	}
	req.Symbol = normalizedSymbol 

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		log.Printf("ERROR: Invalid startDate '%s': %v", req.StartDate, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Must be RFC3339 (e.g., 2025-10-15T00:00:00Z)"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		log.Printf("ERROR: Invalid endDate '%s': %v", req.EndDate, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Must be RFC3339 (e.g., 2025-10-16T00:00:00Z)"})
		return
	}

	marketData, err := ctrl.liveSignalService.GetMarketDataBySymbol(c, req.Symbol)
	if err != nil {
		log.Printf("Warning: Could not fetch market data for symbol %s during signal creation: %v. CurrentPrice will be initialized to 0.", req.Symbol, err)
		req.CurrentPrice = 0 
	} else if marketData != nil {
		req.CurrentPrice = marketData.CurrentPrice
		log.Printf("Fetched live current price for %s: %.4f (overriding client's %.4f)", req.Symbol, req.CurrentPrice, req.CurrentPrice)
	} else {
		log.Printf("No market data found for %s during signal creation. CurrentPrice will be initialized to 0.", req.Symbol)
		req.CurrentPrice = 0
	}

	createdByRole := "Admin"
	creatorID := uint(1) 

	signal := models.Signal{
		TraderID:       creatorID, 
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
		CreatedBy:      createdByRole,
		CreatorID:      creatorID,
	}

	log.Printf("Parsed Signal Data before service call: %+v", signal)

	createdSignal, err := ctrl.liveSignalService.CreateSignal(c, &signal)
	if err != nil {
		log.Printf("ERROR: Failed to create signal in service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create signal"})
		return
	}

	log.Printf("Signal created successfully: ID=%d, Symbol=%s", createdSignal.ID, createdSignal.Symbol)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Signal created successfully",
		"data":    createdSignal,
	})
}

func GetMarketDataAPI(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		log.Println("ERROR: Database connection not found in Gin context for GetMarketDataAPI")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		log.Println("ERROR: Database context object is not a GORM DB instance for GetMarketDataAPI")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database context object is invalid"})
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
			Volume24H:      md.Volume24H,
		})
	}

	c.JSON(http.StatusOK, apiResponse)
}
