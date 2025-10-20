package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log" // Import log package
	"net/http"
	"strconv" // Use strconv.ParseFloat for better error handling
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
		log.Printf("Error fetching market data from %s: %v", url, err)
		return nil, fmt.Errorf("failed to fetch market data from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("External API returned non-OK status: %d %s", resp.StatusCode, resp.Status)
		bodyBytes, readErr := ioutil.ReadAll(resp.Body)
		if readErr == nil {
			errorMessage += fmt.Sprintf(", Response Body: %s", string(bodyBytes))
		}
		log.Println(errorMessage)
		return nil, fmt.Errorf("failed to fetch market data, status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %v", url, err)
		return nil, fmt.Errorf("failed to read market data response: %w", err)
	}

	// Log the raw response body for debugging (truncate if very long)
	log.Printf("Raw API response (first %d chars): %s", min(len(body), 500), body[:min(len(body), 500)])

	var kucoinResponse struct {
		Data struct {
			Ticker []struct {
				Symbol     string `json:"symbol"`
				LastPrice  string `json:"last"`
				ChangeRate string `json:"changeRate"` // 24h change rate
				VolValue   string `json:"volValue"`   // 24h volume
			} `json:"ticker"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &kucoinResponse); err != nil {
		log.Printf("Error unmarshalling market data JSON: %v. Raw body: %s", err, string(body))
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
		"SOL-USDT":  "Solana", // Added for more variety
		"DOGE-USDT": "Dogecoin",
	}

	for _, ticker := range kucoinResponse.Data.Ticker {
		name, ok := symbolsOfInterest[ticker.Symbol]
		if !ok {
			continue // Skip symbols not in our interest list
		}

		lastPrice, pErr := parseFloatStrict(ticker.LastPrice)
		changeRate, cErr := parseFloatStrict(ticker.ChangeRate)
		volume, vErr := parseFloatStrict(ticker.VolValue)

		if pErr != nil {
			log.Printf("Skipping %s due to price parsing error: %v (value: '%s')", ticker.Symbol, pErr, ticker.LastPrice)
			continue
		}
		if cErr != nil {
			log.Printf("Skipping %s due to changeRate parsing error: %v (value: '%s')", ticker.Symbol, cErr, ticker.ChangeRate)
			continue
		}
		if vErr != nil {
			log.Printf("Skipping %s due to volume parsing error: %v (value: '%s')", ticker.Symbol, vErr, ticker.VolValue)
			continue
		}

		marketDataList = append(marketDataList, models.MarketDataAPIResponse{
			Symbol:         ticker.Symbol,
			Name:           name,
			CurrentPrice:   lastPrice,
			PriceChange24H: changeRate * 100, // Convert rate (e.g., 0.0123) to percentage (1.23)
			Volume24H:      volume,
			// Construct a generic logo URL. You might need a more robust service for actual logos.
			LogoURL: fmt.Sprintf("https://cryptoicons.org/api/icon/%s/24", getBaseSymbol(ticker.Symbol)),
		})

		if len(marketDataList) >= 7 { // Limit the number of assets displayed on dashboard
			break
		}
	}

	if len(marketDataList) == 0 {
		log.Println("No relevant live data was successfully processed from the API. Falling back to mock data.")
		// Return mock data ONLY if live data fetching and parsing entirely failed for all symbols
		return s.getMockMarketData(), nil
	}

	log.Printf("Successfully processed %d market data entries.", len(marketDataList))
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
	parts := split(symbol, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return symbol
}

// Simple split function (could use strings.Split)
func split(s, sep string) []string {
	var result []string
	for {
		idx := -1
		for i := 0; i <= len(s)-len(sep); i++ {
			if s[i:i+len(sep)] == sep {
				idx = i
				break
			}
		}
		if idx == -1 {
			result = append(result, s)
			break
		}
		result = append(result, s[:idx])
		s = s[idx+len(sep):]
	}
	return result
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
	return []models.MarketDataAPIResponse{
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
	}
}
