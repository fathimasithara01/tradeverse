package graph

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/graphql-go/graphql"
)

var mockDB = struct {
	UserData      map[string]interface{}
	DashboardData map[string]interface{}
	SignalsData   []map[string]interface{}
	PortfolioData map[string]interface{}
	SettingsData  map[string]interface{}
}{
	UserData:      map[string]interface{}{"id": "user-123", "name": "Sithara", "email": "customer@tradeverse.com", "createdAt": "2024-01-01T10:00:00Z"},
	DashboardData: map[string]interface{}{"metrics": map[string]interface{}{"totalPnl": 1250.75, "winRate": 72.0}, "equityCurve": []map[string]interface{}{{"date": "Jan", "value": 10000}, {"date": "Feb", "value": 10200}, {"date": "May", "value": 11250}}},
	SignalsData: []map[string]interface{}{
		{"id": "signal-1", "traderName": "TraderTwo", "status": "ACTIVE", "entryPrice": 45000.0, "stopLoss": 43000.0, "targetPrice": 50000.0, "currentPrice": 44980.0, "strategy": "Day Trading", "riskLevel": "HIGH", "tradePeriodStart": "5/15/2025", "tradePeriodEnd": "6/14/2025", "createdAt": time.Now().Add(-4 * 24 * time.Hour).Format(time.RFC3339)},
		{"id": "signal-2", "traderName": "TraderThree", "status": "ACTIVE", "entryPrice": 10000.0, "stopLoss": 9500.0, "targetPrice": 12000.0, "currentPrice": 10006.184, "strategy": "Scalping", "riskLevel": "LOW", "tradePeriodStart": "5/10/2026", "tradePeriodEnd": "6/9/2025", "createdAt": time.Now().Add(-482400 * time.Minute).Format(time.RFC3339)},
		{"id": "signal-3", "traderName": "InactiveTrader", "status": "NO_STATUS"},
	},
	PortfolioData: map[string]interface{}{
		"openPositions":   []map[string]interface{}{{"id": "p1", "symbol": "EUR/USD", "type": "BUY", "size": 0.1, "entryPrice": 1.0850, "pnl": 35.25}},
		"closedPositions": []map[string]interface{}{{"id": "p2", "symbol": "GBP/USD", "type": "BUY", "size": 0.2, "closeTime": "2024-05-10 14:30", "pnl": 112.50}},
	},
	SettingsData: map[string]interface{}{
		"brokerConnection": map[string]interface{}{"accountId": "123456", "isConnected": true},
		"security":         map[string]interface{}{"twoFactorEnabled": true},
		"billing":          map[string]interface{}{"planName": "Pro Plan", "nextBillingDate": "July 1, 2025"},
	},
}

func MeResolver(p graphql.ResolveParams) (interface{}, error) { return mockDB.UserData, nil }
func DashboardResolver(p graphql.ResolveParams) (interface{}, error) {
	return mockDB.DashboardData, nil
}
func SignalsResolver(p graphql.ResolveParams) (interface{}, error) { return mockDB.SignalsData, nil }
func PortfolioResolver(p graphql.ResolveParams) (interface{}, error) {
	return mockDB.PortfolioData, nil
}
func SettingsResolver(p graphql.ResolveParams) (interface{}, error) { return mockDB.SettingsData, nil }

func PortfolioUpdatesResolver(p graphql.ResolveParams) (interface{}, error) {
	c := make(chan interface{})
	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if len(mockDB.PortfolioData["openPositions"].([]map[string]interface{})) > 0 {
					pnl := mockDB.PortfolioData["openPositions"].([]map[string]interface{})[0]["pnl"].(float64)
					mockDB.PortfolioData["openPositions"].([]map[string]interface{})[0]["pnl"] = pnl + (rand.Float64()-0.5)*5
				}
				fmt.Println("[REAL-TIME] Pushing portfolio update to client...")
				c <- mockDB.PortfolioData
			case <-ctx.Done():
				fmt.Println("[REAL-TIME] Client disconnected from portfolio updates.")
				close(c)
				return
			}
		}
	}(p.Context)
	return c, nil
}
