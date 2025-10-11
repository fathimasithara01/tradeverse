// internal/admin/cron/cron.go
package cron

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cronn "github.com/robfig/cron/v3"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	customerService "github.com/fathimasithara01/tradeverse/internal/customer/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type CoinGeckoCoin struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	Image                    string  `json:"image"`
	CurrentPrice             float64 `json:"current_price"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}

func FetchAndSaveMarketData(db *gorm.DB) {
	apiURL := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1&sparkline=false&price_change_percentage=24h"

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error fetching market data from API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("API returned non-OK status: %d, Response: %s", resp.StatusCode, string(body))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response body: %v", err)
		return
	}

	var coins []CoinGeckoCoin
	if err := json.Unmarshal(body, &coins); err != nil {
		log.Printf("Error unmarshaling API response: %v", err)
		return
	}

	for _, coin := range coins {
		marketData := models.MarketData{
			Symbol:         strings.ToUpper(coin.Symbol),
			Name:           coin.Name,
			CurrentPrice:   coin.CurrentPrice,
			PriceChange24H: coin.PriceChangePercentage24h,
			LogoURL:        coin.Image,
			// LastUpdated:    time.Now(),
		}
		log.Printf(" current data for %s", coin.CurrentPrice)

		result := db.Where(models.MarketData{Symbol: marketData.Symbol}).Assign(marketData).FirstOrCreate(&marketData)
		if result.Error != nil {
			log.Printf("Error saving/updating market data for %s: %v", coin.Symbol, result.Error)
		} else if result.RowsAffected == 0 {

		} else {
			log.Printf("Saved/Updated market data for %s (Current Price: %.4f)", coin.Symbol, coin.CurrentPrice)
		}
	}
	log.Println("Market data fetch complete.")
}

func StartCronJobs(
	subscriptionService service.ISubscriptionService,
	customerServiceForTraderSubs customerService.AdminSubscriptionService,
	liveSignalService service.ILiveSignalService,
	db *gorm.DB,
) {
	c := cronn.New()

	c.AddFunc("@daily", func() {
		log.Println("Running daily subscription check...")
		if err := subscriptionService.DeactivateExpiredSubscriptions(); err != nil {
			log.Printf("Error checking expired subscriptions: %v", err)
		}
	})

	c.AddFunc("@every 1h", func() {
		log.Println("Running hourly trader subscription status update...")

	})

	c.AddFunc("@every 5m", func() {
		log.Println("Starting market data fetch...")
		FetchAndSaveMarketData(db)
	})

	c.AddFunc("@every 1m", func() {
		log.Println("Starting signal current price update from market data...")
		if err := liveSignalService.UpdateAllSignalsCurrentPrices(context.Background()); err != nil {
			log.Printf("Error updating signal current prices: %v", err)
		}
	})

	c.AddFunc("@every 30s", func() {
		log.Println("Starting signal status evaluation (SL/Target/Activation)...")
		if err := liveSignalService.CheckAndSetSignalStatuses(context.Background()); err != nil {
			log.Printf("Error checking and setting signal statuses: %v", err)
		}
	})

	c.Start()
	log.Println("Cron jobs started.")
}
