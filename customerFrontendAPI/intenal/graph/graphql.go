package graph

import (
	"log"

	"github.com/graphql-go/graphql"
)

var Schema graphql.Schema

func InitSchema() {
	userType := graphql.NewObject(graphql.ObjectConfig{Name: "User", Fields: graphql.Fields{"id": &graphql.Field{Type: graphql.ID}, "name": &graphql.Field{Type: graphql.String}, "email": &graphql.Field{Type: graphql.String}, "createdAt": &graphql.Field{Type: graphql.String}}})
	performanceMetricsType := graphql.NewObject(graphql.ObjectConfig{Name: "PerformanceMetrics", Fields: graphql.Fields{"totalPnl": &graphql.Field{Type: graphql.Float}, "winRate": &graphql.Field{Type: graphql.Float}}})
	equityPointType := graphql.NewObject(graphql.ObjectConfig{Name: "EquityPoint", Fields: graphql.Fields{"date": &graphql.Field{Type: graphql.String}, "value": &graphql.Field{Type: graphql.Float}}})
	assetAllocationType := graphql.NewObject(graphql.ObjectConfig{Name: "AssetAllocation", Fields: graphql.Fields{"assetClass": &graphql.Field{Type: graphql.String}, "percentage": &graphql.Field{Type: graphql.Float}}})
	dashboardType := graphql.NewObject(graphql.ObjectConfig{Name: "Dashboard", Fields: graphql.Fields{"metrics": &graphql.Field{Type: performanceMetricsType}, "equityCurve": &graphql.Field{Type: graphql.NewList(equityPointType)}, "assetAllocation": &graphql.Field{Type: graphql.NewList(assetAllocationType)}}})
	tradeSignalType := graphql.NewObject(graphql.ObjectConfig{Name: "TradeSignal", Fields: graphql.Fields{"id": &graphql.Field{Type: graphql.ID}, "traderName": &graphql.Field{Type: graphql.String}, "status": &graphql.Field{Type: graphql.String}, "entryPrice": &graphql.Field{Type: graphql.Float}, "stopLoss": &graphql.Field{Type: graphql.Float}, "targetPrice": &graphql.Field{Type: graphql.Float}, "currentPrice": &graphql.Field{Type: graphql.Float}, "strategy": &graphql.Field{Type: graphql.String}, "riskLevel": &graphql.Field{Type: graphql.String}, "tradePeriodStart": &graphql.Field{Type: graphql.String}, "tradePeriodEnd": &graphql.Field{Type: graphql.String}, "createdAt": &graphql.Field{Type: graphql.String}}})
	tradeType := graphql.NewObject(graphql.ObjectConfig{Name: "Trade", Fields: graphql.Fields{"id": &graphql.Field{Type: graphql.ID}, "symbol": &graphql.Field{Type: graphql.String}, "type": &graphql.Field{Type: graphql.String}, "size": &graphql.Field{Type: graphql.Float}, "entryPrice": &graphql.Field{Type: graphql.Float}, "pnl": &graphql.Field{Type: graphql.Float}, "closeTime": &graphql.Field{Type: graphql.String}}})
	monthlyPnlType := graphql.NewObject(graphql.ObjectConfig{Name: "MonthlyPnl", Fields: graphql.Fields{"month": &graphql.Field{Type: graphql.String}, "pnl": &graphql.Field{Type: graphql.Float}}})
	portfolioAnalyticsType := graphql.NewObject(graphql.ObjectConfig{Name: "PortfolioAnalytics", Fields: graphql.Fields{"monthlyPerformance": &graphql.Field{Type: graphql.NewList(monthlyPnlType)}}})
	portfolioType := graphql.NewObject(graphql.ObjectConfig{Name: "Portfolio", Fields: graphql.Fields{"openPositions": &graphql.Field{Type: graphql.NewList(tradeType)}, "closedPositions": &graphql.Field{Type: graphql.NewList(tradeType)}, "analytics": &graphql.Field{Type: portfolioAnalyticsType}}})
	brokerConnectionType := graphql.NewObject(graphql.ObjectConfig{Name: "BrokerConnection", Fields: graphql.Fields{"accountId": &graphql.Field{Type: graphql.String}, "isConnected": &graphql.Field{Type: graphql.Boolean}}})
	userSecurityType := graphql.NewObject(graphql.ObjectConfig{Name: "UserSecurity", Fields: graphql.Fields{"twoFactorEnabled": &graphql.Field{Type: graphql.Boolean}}})
	billingType := graphql.NewObject(graphql.ObjectConfig{Name: "Billing", Fields: graphql.Fields{"planName": &graphql.Field{Type: graphql.String}, "nextBillingDate": &graphql.Field{Type: graphql.String}}})
	settingsType := graphql.NewObject(graphql.ObjectConfig{Name: "Settings", Fields: graphql.Fields{"brokerConnection": &graphql.Field{Type: brokerConnectionType}, "security": &graphql.Field{Type: userSecurityType}, "billing": &graphql.Field{Type: billingType}}})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{
		"me":        &graphql.Field{Type: userType, Resolve: MeResolver},
		"dashboard": &graphql.Field{Type: dashboardType, Resolve: DashboardResolver},
		"signals":   &graphql.Field{Type: graphql.NewList(tradeSignalType), Resolve: SignalsResolver},
		"portfolio": &graphql.Field{Type: portfolioType, Resolve: PortfolioResolver},
		"settings":  &graphql.Field{Type: settingsType, Resolve: SettingsResolver},
	}})

	rootSubscription := graphql.NewObject(graphql.ObjectConfig{Name: "Subscription", Fields: graphql.Fields{
		"portfolioUpdates": &graphql.Field{Type: portfolioType, Resolve: PortfolioUpdatesResolver},
	}})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery, Subscription: rootSubscription})
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}
}
