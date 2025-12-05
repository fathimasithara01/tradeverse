
A production-ready multi-role fintech platform with real-time trading,
wallet ledgering, subscription automation, and admin controls — built using
Golang, PostgreSQL, Clean Architecture, and Cron Workers.

# Tradeverse – Multi-Role Trading Automation Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)]()
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Build Status](https://img.shields.io/badge/Status-Production--Ready-success)]()
[![PRs Welcome](https://img.shields.io/badge/PRs-Welcome-blue.svg)]()

A production-grade **trading, wallet, and subscription automation system** built using
**Golang, PostgreSQL, Clean Architecture, Cron Jobs, RBAC, and Real-Time Market APIs.**

Tradeverse is designed as a **multi-role fintech platform** supporting:

* **Admin**
* **Trader**
* **Customer**

It follows strict engineering standards similar to real-world fintech systems.

---

##  Overview

Tradeverse includes:

* Real-time market price fetching (OHLC / snapshot)
* Automated subscription lifecycle (expiry, disable, notify)
* Wallet system with strict ledger consistency
* Trader signal publishing + analytics
* Dynamic admin-driven pricing
* JWT + RBAC (Admin / Trader / Customer)
* Cron-based background workers
* Full Admin UI

This project demonstrates **scalable backend architecture**, **clean code**, and **production-level fintech engineering**.

---

##  Features

###  Cron Jobs (Schedulers)

* Real-time market data fetcher
* Configurable intervals
* Retry + rate limiting
* Normalized OHLC storage

###  Subscription Automation

* Auto-expiry
* Auto-disable access
* Notification triggers
* Batch-optimized

###  Admin UI

* User CRUD
* Commission & pricing configuration
* API key / cron intervals
* Dashboard + charts

###  Dynamic Pricing

* Commission % control
* DB-driven pricing
* Trader-specific plans

###  Signal Cards

* Entry, SL, targets
* Live price
* Color-coded cards
* Trader info

###  RBAC (Role Based Access Control)

* Admin / Trader / Customer
* JWT with role + plan + expiry
* Backend + UI-controlled access

###  User Management

* CRUD
* Role assignment
* Block/Unblock
* Subscription status

###  Dashboard Insights

* Revenue analytics
* Live prices
* Active subscriptions
* Trader/customer statistics

---

##  Role Breakdown

###  Customer

* Signup / Login
* Browse traders
* Subscribe / Unsubscribe
* Wallet deposit/withdraw/history
* KYC uploads
* View subscribed signals

###  Trader

* Create/publish signals
* Live trade updates
* Subscription plans
* Subscriber analytics
* Profile management

###  Admin

* Manage users & traders
* Pricing & commission rules
* System settings
* Dashboard analytics
* Audit logs

---

##  Architecture Overview

###  Clean Architecture Layers

* **Handlers** — routing, request parsing
* **Services** — business logic
* **Repositories** — DB + external APIs
* **Domain Models** — entities

###  Project Structure

```
tradeverse/
├── cmd/
│   ├── admin/
│   ├── trader/
│   └── customer/
├── config/
├── internal/
│   ├── admin/
│   ├── trader/
│   └── customer/
├── migrations/
├── pkg/
│   ├── auth/
│   ├── models/
│   ├── payment_gateway.go
│   ├── seeder/
│   └── utils/
├── static/
├── templates/
└── README.md
```

###  Request Flow

```
Client → Router → Handler → Service → Repository → PostgreSQL
```

---

##  Core Modules

###  Authentication

* JWT-based
* RBAC middleware
* Claims contain role + subscription

###  Wallet System

* Deposit / Withdraw
* Transaction ledger
* Race-condition safe
* Accurate balance tracking

###  Trader Module

* Signal CRUD
* Live trade publishing
* Subscription plans
* Subscriber analytics

###  Customer Module

* Explore traders
* Subscribe/unsubscribe
* View signals
* Wallet + KYC

###  Admin Module

* User & trader management
* Plans & commissions
* System configuration
* Analytics dashboard

---

##  API Endpoints (High-Level)

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

## 🛠 Tech Stack

* **Golang (Gin Framework)**
* **PostgreSQL**
* **GORM ORM**
* **Cron Jobs**
* **Server-rendered Admin UI**
* **JWT Authentication**
* **Clean Architecture + DDD**
* **Docker Ready**

---

##  Internals Summary

* Market Fetcher — cron-based OHLC + price
* Subscription Watcher — auto-expiry
* RBAC Engine — role logic
* Admin Configuration
* Signal Cards with real-time data
* Dashboard analytics

---

##  Running the Project

### Start Admin

```sh
go run cmd/admin/main.go
```

### Start Trader

```sh
go run cmd/trader/main.go
```

### Start Customer

```sh
go run cmd/customer/main.go
```

### Run Migrations

```sh
go run internal/migrations/main.go
```

### Seed Data

```sh
go run pkg/seeder/main.go
```

---

##  Run Locally

```bash
git clone https://github.com/fathimasithara01/tradeverse
cd tradeverse
cp .env.example .env
go mod tidy
go run cmd/server/main.go


##  Security

* JWT expiry + rotation
* Secrets via Vault / AWS SM
* SQL injection protection
* HTTPS + Nginx
* CORS restrictions
* Rate limiting

---

##  Deployment

* Docker / Docker Compose
* Kubernetes-ready
* Worker containers
* Prometheus metrics
* Redis caching & pub/sub
* Cloud PostgreSQL

---

##  Why This Project is Valuable for Recruiters
- Shows backend ownership end-to-end
- Demonstrates realistic fintech domain knowledge
- Highlights distributed cron workers & automation
- Proves understanding of architecture & scalability
- Perfect fit for Backend (Golang) / Fintech / SaaS roles


##  Why This Project Stands Out

* Rare multi-role fintech system
* Realistic wallet + subscription engine
* Clean Architecture + DDD
* Horizontally scalable
* Production-like engineering

---


 Postman Collection (API Testing)

A complete Postman Collection is included to help you test all TradeVerse APIs easily.

 What’s Included

The collection covers:

Authentication (Admin, Trader, Customer)

User Profile & Management

Wallet (Deposit, Withdraw, Transactions)

Payments & Subscription Lifecycle

Trader Signals

Copy Trading Automation

KYC Verification

Cron Job Simulation APIs

Download Postman Collection

You can find the collection file inside the repository:

/postman/TradeVerse_API_Collection.json

 How to Import & Use

Open Postman

Click Import

Select:
postman/TradeVerse_API_Collection.json

Set the following Environment Variables:

BASE_URL = http://localhost:8080
ADMIN_TOKEN = <set after admin login>
TRADER_TOKEN = <set after trader login>
CUSTOMER_TOKEN = <set after customer login>


Start testing APIs 

 Why This Section Is Important

Shows API completeness

HR & Interviewers can easily test your backend

Makes the project look like a real production SaaS

Increases credibility and professional quality

##  System Diagram

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

##  Author

**Fathima Sithara**
Backend Developer — Golang • Microservices • Fintech
GitHub: [https://github.com/fathimasithara01](https://github.com/fathimasithara01)

---
