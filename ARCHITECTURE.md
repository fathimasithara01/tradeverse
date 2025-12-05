# TradeVerse – Architecture Overview

This document provides a professional, high‑level architectural overview of the **TradeVerse Multi‑Role Trading Automation Platform**, covering system design principles, service structure, internal layers, module responsibilities, and data flows.

---

# 1. System Architecture Summary

TradeVerse is designed using **Clean Architecture** and **Domain‑Driven Design (DDD)** to ensure:

* Strong separation of concerns
* Scalability across Admin, Trader, and Customer modules
* Maintainability and testability
* Production‑grade modular design
* Clear domain boundaries and use‑case isolation

### **Core Request Flow**

```
Client → Router → Handlers → Services (Use‑Cases) → Repository → PostgreSQL
```

All external interfaces depend on abstractions, not implementations.

---

# 2. Service Execution Model (cmd/)

The platform runs **three independent services**, each acting as a separate executable. This makes the system horizontally scalable and role‑specialized.

### **Admin Service – `cmd/admin/`**

* System management
* User lifecycle control
* Pricing, commissions, analytics

### **Trader Service – `cmd/trader/`**

* Signal publishing
* Subscription plan creation
* Live trade broadcasting

### **Customer Service – `cmd/customer/`**

* Wallet operations
* Subscriptions
* KYC verification

Each service initializes:

* Router
* Config loader
* Middlewares
* Dependency injection
* Cron jobs (where required)

---

# 3. Clean Architecture Layers

### **3.1 Handlers (Transport Layer)**

* HTTP request parsing & validation
* Invokes service layer
* Sends standardized JSON responses

### **3.2 Services / Use‑Cases (Business Logic)**

Implements core domain rules, including:

* Wallet consistency
* Subscription lifecycle
* RBAC validation
* Signal creation & publishing

> This layer contains the primary business logic and is framework‑agnostic.

### **3.3 Repository Layer (Persistence)**

* DB operations using GORM
* Transaction‑safe debits/credits
* Optimized queries with indexing

### **3.4 Domain Models**

* Pure Go structs representing business entities

Examples:

* User
* Wallet & Transactions
* Subscription
* Signal
* MarketPrice

---

# 4. Folder Structure Overview

```
tradeverse/
├── cmd/
│   ├── admin/                 # Admin service main entry
│   ├── trader/                # Trader service main entry
│   └── customer/              # Customer service main entry
│
├── internal/
│   ├── admin/                 # Admin-specific modules
│   ├── trader/                # Trader logic
│   ├── customer/              # Customer workflows
│   └── migrations/            # DB migrations
│
├── pkg/
│   ├── auth/                  # JWT & RBAC utilities
│   ├── models/                # Domain entities
│   ├── seeder/                # Initial data seeding
│   └── utils/                 # Common helpers
│
├── config/                    # Configuration files
├── static/                    # Admin UI assets
├── templates/                 # Server-rendered Admin UI
└── README.md
```

---

# 5. Module Responsibilities

### **Admin Module**

* User & trader management
* Pricing & commissions
* System settings
* Dashboard analytics

### **Trader Module**

* Create/update signals
* Publish live trades
* Manage subscription plans
* Monitor subscribers

### **Customer Module**

* Wallet operations
* Subscriptions
* KYC processing
* Fetch subscribed signals

### **Migrations Module**

* Schema definitions
* Automated migrations

---

# 6. Background Workers

### **6.1 Market Price Fetcher**

* Scheduled polling
* Data normalization
* OHLC storage

### **6.2 Subscription Lifecycle Manager**

* Validates expired plans
* Deactivates access
* Pushes notifications

---

# 7. Database Overview

### **Primary Tables**

* users
* wallets
* wallet_transactions
* subscription_plans
* subscriptions
* trader_signals
* market_prices
* kyc_documents
* settings

### **DB Guarantees**

* Strong referential integrity
* Indexed critical columns
* Wallet operations use ACID transactions

---

# 8. High-Level Workflow Examples

### **Customer Views Signals**

```
Customer → Login → Subscribe to Trader → Signal Published → Customer Fetches Subscribed Signals
```

### **Real‑Time Price Scheduler**

```
Scheduler → Fetch Market API → Normalize Data → Store in DB → Dashboard Uses Data
```

---

# 9. Summary

This architecture ensures a **robust, maintainable, and production‑ready trading platform**, with clear module separation, clean domain boundaries, and scalable service design. It reflects industry best practices suitable for fintech‑grade applications.
