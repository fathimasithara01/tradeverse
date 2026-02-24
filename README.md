# TradeVerse — Multi-Role Trading Platform (Golang + PostgreSQL)

TradeVerse is a backend trading platform built using Go (Gin), PostgreSQL, and Clean Architecture principles.  
It supports multiple user roles (Admin, Trader, Customer) and includes wallet management, subscription automation, and signal publishing.

TradeVerse is a modular monolith backend application built using Go (Gin) and PostgreSQL. It implements role-based trading workflows, wallet accounting, and subscription lifecycle management using structured service-layer business logic.
---

##  Overview

TradeVerse is designed as a **modular monolith** with clear separation of concerns using Clean Architecture.

It includes:

- Role-based authentication (Admin / Trader / Customer)
- Wallet system with transaction ledger
- Subscription lifecycle automation
- Trader signal publishing
- Admin configuration and pricing management
- Cron-based background jobs

---

##  Tech Stack

- Go (Gin Framework)
- PostgreSQL
- GORM
- JWT Authentication
- Cron Jobs (Schedulers)
- Server-rendered Admin UI
- Clean Architecture (Handler → Service → Repository)
- Docker (local setup)

---

##  Architecture

The project follows layered architecture with clear separation:

- Handler Layer — HTTP request parsing and routing
- Service Layer — Business rules and transaction orchestration
- Repository Layer — Database access using GORM
- Domain Models — Core entities and validation logic

### Request Flow

Client → Router → Handler → Service → Repository → PostgreSQL

---

## Roles & Capabilities

### Customer
- Signup / Login  
- Browse traders  
- Subscribe / Unsubscribe  
- Deposit / Withdraw wallet funds  
- View transaction history  
- View subscribed trader signals  

### Trader
- Create and manage trading signals  
- Publish live trades  
- Create subscription plans  
- View subscriber information  

### Admin
- Manage users and traders  
- Configure pricing and commissions  
- Monitor subscriptions  
- View basic analytics  

---

## Wallet System

- Deposit and withdraw functionality  
- Transaction ledger stored in PostgreSQL  
- Balance updates handled through service layer logic  
- Balance updates are executed within database transactions to ensure atomicity and prevent inconsistent wallet states.
---

## ⏱ Subscription Automation

Subscription status is validated using scheduled cron jobs:

- Automatic expiry of inactive subscriptions
- Access restriction after expiry
- Periodic verification of subscription validity
---

## 📈 Market Data

- Periodic fetching of market price data (OHLC format)  
- Stored in normalized database tables  
- Used for trader signal context  

---

## Project Structure

tradeverse/
├── cmd/
│ ├── admin/
│ ├── trader/
│ └── customer/
├── config/
├── internal/
│ ├── admin/
│ ├── trader/
│ └── customer/
├── migrations/
├── pkg/
│ ├── auth/
│ ├── models/
│ ├── seeder/
│ └── utils/
├── static/
├── templates/
└── README.md

---

## 🛠 Running Locally

### 1️ Clone Repository

git clone https://github.com/fathimasithara01/tradeverse
cd tradeverse

### 2️ Setup Environment

cp .env.example .env

go mod tidy

### 3️ Run Migrations

go run internal/migrations/main.go

### 4️ Seed Data (Optional)

go run pkg/seeder/main.go

### 5️ Start Application

go run cmd/server/main.go

---

 ## What This Project Demonstrates

- Backend ownership of a complete domain

- Role-based access control (RBAC)

- Clean Architecture implementation

- Structured business logic separation

- Database transaction handling
 
- Cron-based background processing

## Limitations

- Deployment: Local Docker setup for development and testing.
  
- No distributed scaling setup

- No Kubernetes or cloud deployment included

- Intended for backend system design demonstration
  
---

## Author

Fathima Sithara
Backend Engineer (Golang)
GitHub: https://github.com/fathimasithara01
