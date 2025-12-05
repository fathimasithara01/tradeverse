#  TradeVerse – Multi-Role Copy Trading Platform

A Production-Ready SaaS Application — **Admin, Trader & Customer Modules**

TradeVerse is a complete **Copy Trading SaaS Platform** built using **Golang**, following **Clean Architecture**, **Domain-Driven Design (DDD)**, and a **microservice-ready modular structure**. It supports **Admin**, **Trader**, and **Customer** roles with secure authentication, wallet operations, subscriptions, signal publishing, KYC, and performance analytics.

---

##  Key Highlights

* ✔ **Clean Architecture + DDD**
* ✔ **Role-based modular services** (`cmd/admin`, `cmd/trader`, `cmd/customer`)
* ✔ **Secure JWT Authentication**
* ✔ **Fully validated Wallet System**
* ✔ **Subscriptions + Live Trading Signals**
* ✔ **PostgreSQL + GORM ORM**
* ✔ **Scalable & Microservice-Ready Structure**

---

##  Role Overview

### ** Customer**

* Signup/Login (JWT)
* Browse traders
* Subscribe/unsubscribe
* KYC upload & status tracking
* View signals from subscribed traders
* Wallet: deposit, withdrawal, transaction history

### ** Trader**

* Create trading signals
* Publish live trades
* Create/manage subscription plans
* View subscribers
* Manage trader profile

### ** Admin**

* Manage traders & system data
* Manage subscription plans
* Monitor activity & audits

---

##  Architecture Overview

### **Clean Architecture Layers**

* **Handlers / Controllers** – Request validation + routing
* **Services / Use-Cases** – Core business logic
* **Repositories** – Database interactions using GORM
* **Domain Models** – Independent business entities

### **Modular Directory Structure**

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
├── migrations/
│
├── pkg/
│   ├── auth/
│   ├── models/
│   ├── payment_gateway.go
│   ├── seeder/
│   └── utils/
│
├── static/
├── templates/
└── README.md

### ** Data Flow**

```
Client → Router → Handler → Service → Repository → PostgreSQL
```

---

##  Core Features

### **Authentication & Access Control**

* JWT login/signup
* Role-based access (Admin / Trader / Customer)
* Token validation middleware
* Session management

###  **Wallet System**

* Deposit / Withdraw
* Balance compute
* Transaction history
* Multi-role actions
* Strong validations to prevent corruption

###  **Trader Module**

* CRUD Trading Signals
* Live trade publishing
* Subscription plan creation
* Subscriber list view
* Profile CRUD

###  **Customer Module**

* View traders & performance metrics
* Subscribe/unsubscribe
* View signals from subscribed traders
* KYC upload & verification
* Full profile management

###  **Admin Module**

* Manage traders & customers
* Manage subscription plans
* Audit & reporting utilities

---

##  API Overview (High-Level)

###  **Trader**

* `/login`
* `/createSignal`, `/getAllSignals`, `/updateSignal`
* `/CreateTraderSubscriptionPlan`, `/ListSubscribers`
* `/PublishLiveTrade`
* `/GetBalance`, `/Deposit`, `/Withdraw`

###  **Customer**

* `/signup`, `/login`
* `/ListTraders`, `/GetTraderDetails`
* `/SubscribeToTrader`, `/GetSignalsFromSubscribedTraders`
* `/kycDocument`, `/GetWalletSummary`

###  **Admin**

* `/ListAdminSubscriptionPlans`
* `/SubscribeToAdminPlan`
* `/CancelAdminSubscription`

> **Full API Documentation available inside the repository.**

---

##  Tech Stack

* **Go (Golang) – Gin / net/http**
* **PostgreSQL**
* **GORM ORM**
* **JWT Authentication**
* **Clean Architecture + DDD**
* **Docker-ready setup**
* **Seeders + Migrations included**

---

##  Running the Services

### **Admin Service**

```
go run cmd/admin/main.go
```

### **Trader Service**

```
go run cmd/trader/main.go
```

### **Customer Service**

```
go run cmd/customer/main.go
```

---

##  Database Migrations

```
go run internal/migrations/main.go
```

##  Database Seeder

```
go run pkg/seeder/main.go
```

---

##  Author

**Fathima Sithara**
Backend Developer (Golang | Microservices)
🔗 **GitHub:** [https://github.com/fathimasithara01](https://github.com/fathimasithara01)

---