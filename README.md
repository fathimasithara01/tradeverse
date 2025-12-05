# Tradeverse – Multi-Role Trading Automation Platform

A production-grade **trading, wallet, and subscription automation system** built using **Go (Golang)**, **PostgreSQL**, **Clean Architecture**, **Cron Jobs**, and **Real-Time Market APIs**.

---

##  Overview
Tradeverse is a complete multi-role fintech platform for **Admin**, **Trader**, and **Customer** workflows.
It includes:
- Real-time market price fetching
- Automated subscription lifecycle
- Advanced RBAC
- Wallet system with strong validations
- Trader signal publishing
- Fully functional Admin UI

This project demonstrates scalable backend engineering, DDD, and production-style fintech workflows.

---

## Implemented Features

### Cron Jobs (Schedulers)
- Periodic real-time market data fetcher
- Configurable intervals via UI
- Retry/backoff + rate-limit safe
- Stores normalized OHLC/price snapshots

### Subscription Automation
- Auto-check for expired plans
- Instantly disables access
- Sends notifications
- Optimized batch processing

### Admin UI
- User management (CRUD, roles, block/unblock)
- Commission & pricing configuration
- System settings (API keys, intervals, toggles)
- Dashboard with charts, analytics & live metrics

### Dynamic Pricing
- Admin sets commission %
- Role-based pricing
- Plan pricing adjustments
- Fully persisted to DB

### Signal Cards
Show: current price, entry, SL, targets, timestamp, trader info.

### RBAC
- Admin / Trader / Customer
- JWT with role + plan
- Server + UI enforcement

### User Management
- CRUD
- Role assignment
- Status & subscription management

### Dashboard
- Traders/customers stats
- Revenue graph
- Active subscriptions
- Live price feed
- Recent signals + charts

---

## Role Overview

### Customer
- Signup/Login
- Browse traders
- Subscribe/unsubscribe
- Upload KYC
- View subscribed signals
- Wallet (deposit/withdraw/history)

### Trader
- Create & publish trading signals
- Push live trades
- Manage subscription plans
- View subscribers
- Profile management

### Admin
- Manage users & traders
- System configuration
- Subscription plans
- Dashboard & analytics
- Audit logs

---

## Architecture Overview

### Clean Architecture
- **Handlers** — routing & validation
- **Services** — core business logic
- **Repositories** — data access
- **Domain Models** — pure business entities

### Project Structure
```
tradeverse/
├── cmd/
│   ├── admin/
│   ├── trader/
│   └── customer/
│
├── config/
│
├── internal/
│   ├── admin/
│   ├── trader/
│   ├── customer/
│
│──── migrations/
│
├── pkg/
│   ├── auth/
│   ├── models/
│   ├── payment_gateway.go/
│   ├── seeder/
│   └── utils/
│
├── static/
├── templates/
└── README.md
```

### Request Flow
```
Client → Router → Handler → Service → Repository → PostgreSQL
```

---

## Core Modules

### Authentication
- JWT-based
- RBAC middleware
- Claims store role + subscription info

### Wallet System
- Deposit / Withdraw
- Transaction history
- Race-condition safe
- Ledger accuracy guaranteed

### Trader Module
- CRUD trading signals
- Publish live trades
- Subscription plans
- Subscriber management

### Customer Module
- Explore traders
- Subscribe/unsubscribe
- View signals
- Upload KYC
- Wallet summary

### Admin Module
- Manage users/traders/customers
- Manage subscription plans
- Dashboard & analytics

---

## API Endpoints (High-Level)

### Trader
```
/login
/createSignal
/getAllSignals
/updateSignal
/CreateTraderSubscriptionPlan
/ListSubscribers
/PublishLiveTrade
/GetBalance
/Deposit
/Withdraw
```

### Customer
```
/signup
/login
/ListTraders
/GetTraderDetails
/SubscribeToTrader
/GetSignalsFromSubscribedTraders
/kycDocument
/GetWalletSummary
```

### Admin
```
/ListAdminSubscriptionPlans
/SubscribeToAdminPlan
/CancelAdminSubscription
```

---

## Tech Stack
- Go (Golang) — Gin
- PostgreSQL
- GORM ORM
- Cron Jobs
- Server-side Admin UI
- JWT Authentication
- Clean Architecture + DDD
- Docker-ready

---

## How Internals Work
1. **Market Fetcher** — scheduled, normalized prices stored, live updates.
2. **Subscription Watcher** — expires plans, notifications.
3. **RBAC Engine** — JWT claim inspection.
4. **Admin Panel** — system settings, commission, pricing.
5. **Signal Cards** — color-coded, real-time enriched.
6. **Dashboard** — charts, analytics, revenue, signals.

---

## Running the Project
```
go run cmd/admin/main.go
go run cmd/trader/main.go
go run cmd/customer/main.go
```

### Migrations
```
go run internal/migrations/main.go
```

### Seeder
```
go run pkg/seeder/main.go
```

---

## Security
- JWT expiry & rotation
- External secrets (Vault / AWS Secrets Manager)
- Rate limiting
- SQL injection protection
- HTTPS + Nginx reverse proxy
- Strict CORS

---

## Deployment
- Docker / Docker Compose
- Kubernetes-ready
- Separate worker containers for cron jobs
- Prometheus metrics
- Redis (cache/pub-sub)
- Managed PostgreSQL

---

## Why This Project Stands Out
- Rare multi-role fintech system
- Realistic wallet + subscription engine
- Clean Architecture + DDD
- Horizontally scalable services
- Production-like engineering

---

## System Diagram
```
           +-----------------------+
           |      Client (UI)      |
           +-----------+-----------+
                       |
                       v
              +--------+--------+
              |     API Layer    |
              |   (Gin Handlers) |
              +--------+--------+
                       |
                       v
              +--------+--------+
              |     Services     |
              | (Business Logic) |
              +--------+--------+
                       |
                       v
              +--------+--------+
              |   Repositories   |
              |    (DB Layer)    |
              +--------+--------+
                       |
                       v
              +-------------------+
              |   PostgreSQL DB   |
              +-------------------+
```

---

## Author
**Fathima Sithara**
Backend Developer (Golang • Microservices • Full Stack Capable)
GitHub: https://github.com/fathimasithara01
