package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings" // Added for strings.Split
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type IMarketDataService interface {
	GetLiveMarketData() ([]models.MarketDataAPIResponse, error)
}

type MarketDataService struct {
	// Potentially inject a MarketDataRepository here if you were storing historical data
}

func NewMarketDataService() IMarketDataService {
	return &MarketDataService{}
}

// GetLiveMarketData fetches live market data from the KuCoin API
func (s *MarketDataService) GetLiveMarketData() ([]models.MarketDataAPIResponse, error) {
	const url = "https://api.kucoin.com/api/v1/market/allTickers" // KuCoin's public ticker API

	log.Println("Attempting to fetch market data from:", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("ERROR: Failed to fetch market data from %s: %v", url, err)
		return nil, fmt.Errorf("failed to fetch market data from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("ERROR: Could not read error response body from %s: %v", url, readErr)
		}
		errorMessage := fmt.Sprintf("ERROR: External API returned non-OK status: %d %s, Response Body: %s", resp.StatusCode, resp.Status, string(bodyBytes))
		log.Println(errorMessage)
		return nil, fmt.Errorf("failed to fetch market data, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body from %s: %v", url, err)
		return nil, fmt.Errorf("failed to read market data response: %w", err)
	}

	// Log the raw response body for debugging (truncate if very long)
	log.Printf("DEBUG: Raw API response (first %d chars): %s", min(len(body), 1000), body[:min(len(body), 1000)])

	var kucoinResponse struct {
		Data struct {
			Ticker []struct {
				Symbol     string `json:"symbol"`
				LastPrice  string `json:"last"`
				ChangeRate string `json:"changeRate"` // 24h change rate (e.g., "0.0123")
				VolValue   string `json:"volValue"`   // 24h volume in quote currency (e.g., USDT)
			} `json:"ticker"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &kucoinResponse); err != nil {
		log.Printf("ERROR: Failed to unmarshal market data JSON: %v. Raw body: %s", err, string(body))
		return nil, fmt.Errorf("failed to parse market data from API: %w", err)
	}

	var marketDataList []models.MarketDataAPIResponse
	// Define the specific symbols you are interested in and their display names
	symbolsOfInterest := map[string]string{
		"BTC-USDT":  "Bitcoin",
		"ETH-USDT":  "Ethereum",
		"XRP-USDT":  "Ripple",
		"LTC-USDT":  "Litecoin",
		"ADA-USDT":  "Cardano",
		"SOL-USDT":  "Solana",
		"DOGE-USDT": "Dogecoin",
	}

	for _, ticker := range kucoinResponse.Data.Ticker {
		name, ok := symbolsOfInterest[ticker.Symbol]
		if !ok {
			continue // Skip symbols not in our interest list
		}

		// Parse values with robust error handling
		lastPrice, pErr := parseFloatStrict(ticker.LastPrice)
		if pErr != nil {
			log.Printf("WARN: Skipping %s due to price parsing error: %v (value: '%s')", ticker.Symbol, pErr, ticker.LastPrice)
			continue
		}

		changeRate, cErr := parseFloatStrict(ticker.ChangeRate)
		if cErr != nil {
			log.Printf("WARN: Skipping %s due to changeRate parsing error: %v (value: '%s')", ticker.Symbol, cErr, ticker.ChangeRate)
			continue
		}

		volume, vErr := parseFloatStrict(ticker.VolValue)
		if vErr != nil {
			log.Printf("WARN: Skipping %s due to volume parsing error: %v (value: '%s')", ticker.Symbol, vErr, ticker.VolValue)
			continue
		}

		// Generate Logo URL
		baseSymbol := getBaseSymbol(ticker.Symbol)
		logoURL := fmt.Sprintf("https://cryptoicons.org/api/icon/%s/24", strings.ToLower(baseSymbol)) // Use lowercase for cryptoicons.org
		log.Printf("DEBUG: Generated LogoURL for %s: %s", ticker.Symbol, logoURL)                     // Log generated URL

		marketDataList = append(marketDataList, models.MarketDataAPIResponse{
			Symbol:         ticker.Symbol,
			Name:           name,
			CurrentPrice:   lastPrice,
			PriceChange24H: changeRate * 100, // Convert rate (e.g., 0.0123) to percentage (1.23)
			Volume24H:      volume,
			LogoURL:        logoURL,
		})

		if len(marketDataList) >= 7 { // Limit the number of assets displayed on dashboard
			break
		}
	}

	if len(marketDataList) == 0 {
		log.Println("WARN: No relevant live data was successfully processed from the API. Falling back to mock data.")
		// Return mock data ONLY if live data fetching and parsing entirely failed for all symbols
		return s.getMockMarketData(), nil
	}

	log.Printf("INFO: Successfully processed %d market data entries.", len(marketDataList))
	return marketDataList, nil
}

// parseFloatStrict uses strconv.ParseFloat for more robust error handling
func parseFloatStrict(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("cannot parse empty string to float")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' as float: %w", s, err)
	}
	return f, nil
}

// getBaseSymbol extracts the base currency symbol (e.g., BTC from BTC-USDT)
func getBaseSymbol(symbol string) string {
	parts := strings.Split(symbol, "-") // Using standard library strings.Split
	if len(parts) > 0 {
		return parts[0]
	}
	return symbol
}

// min helper function for logging
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getMockMarketData returns static mock data. Updated with more variety.
func (s *MarketDataService) getMockMarketData() []models.MarketDataAPIResponse {
	// Use a fixed time for consistency if the actual API fails
	now := time.Now()
	mockData := []models.MarketDataAPIResponse{
		{
			Symbol:         "BTC-USDT",
			Name:           "Bitcoin",
			CurrentPrice:   43000.50 + float64(now.Second()%100), // Simulate minor fluctuations
			PriceChange24H: 1.25,
			Volume24H:      25000000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/btc/24",
		},
		{
			Symbol:         "ETH-USDT",
			Name:           "Ethereum",
			CurrentPrice:   2300.75 + float64(now.Second()%50),
			PriceChange24H: -0.78,
			Volume24H:      12000000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/eth/24",
		},
		{
			Symbol:         "XRP-USDT",
			Name:           "Ripple",
			CurrentPrice:   0.58 + float64(now.Second()%10)/100,
			PriceChange24H: 3.10,
			Volume24H:      1500000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/xrp/24",
		},
		{
			Symbol:         "LTC-USDT",
			Name:           "Litecoin",
			CurrentPrice:   70.15 + float64(now.Second()%5)/10,
			PriceChange24H: -1.50,
			Volume24H:      800000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/ltc/24",
		},
		{
			Symbol:         "ADA-USDT",
			Name:           "Cardano",
			CurrentPrice:   0.45 + float64(now.Second()%2)/100,
			PriceChange24H: 0.85,
			Volume24H:      600000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/ada/24",
		},
		{
			Symbol:         "SOL-USDT",
			Name:           "Solana",
			CurrentPrice:   192.00 + float64(now.Second()%20)/10,
			PriceChange24H: 2.10,
			Volume24H:      3000000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/sol/24",
		},
		{
			Symbol:         "DOGE-USDT",
			Name:           "Dogecoin",
			CurrentPrice:   0.20 + float64(now.Second()%5)/1000,
			PriceChange24H: 5.00,
			Volume24H:      2000000000.00,
			LogoURL:        "https://cryptoicons.org/api/icon/doge/24",
		},
	}
	log.Println("INFO: Returning mock market data.")
	return mockData
}
