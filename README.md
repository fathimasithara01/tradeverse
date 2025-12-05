Tradeverse – Multi-Role Trading Automation Platform

A production-grade trading, wallet, and subscription automation system built using Go (Golang), PostgreSQL, Clean Architecture, cron jobs, and real-time market APIs.

🚀 Overview

Tradeverse is a complete multi-role fintech platform designed for Admin, Trader, and Customer workflows.
It supports:

✅ Real-time market price fetching
✅ Automated subscription lifecycle
✅ Advanced RBAC
✅ Wallet with strong validations
✅ Trader signal publishing
✅ Admin UI for configuration and management

This project demonstrates my capabilities in scalable backend engineering, DDD, production architecture, and fintech-grade workflows.

⭐ Implemented Features
🔄 Cron Jobs (Schedulers)

Fetches real-time market data periodically

Configurable intervals via UI

Retry/backoff + rate limit handling

Stores normalized OHLC/price snapshots

Powers dashboards & signal updates

🧾 Subscription Automation

Auto-checks for expired plans

Disables access instantly

Sends notifications

Batch updates for efficiency

🛠️ Admin UI

A functional admin panel with:

User management (CRUD, block/unblock, role assignment)

Commission & pricing management

System configuration (API keys, intervals, feature toggles)

Dashboard with charts, analytics, and live metrics

💰 Commission & Dynamic Pricing

Admin can configure:

Percentage commission

Role-based pricing

Plan pricing

Entire system persists changes in DB

📈 Signal Cards

Each trading signal displays:
Current price, Entry, SL, Targets, Timestamp, and Trader info.

🔐 RBAC (Role-Based Access Control)

Admin / Trader / Customer roles

Enforced both server-side and UI-side

JWT tokens carry role + plan info

👥 User Management

CRUD

Role assignment

Status management

Subscription management

📊 Dashboard

Total traders/customers

Revenue graph

Active subscriptions

Live price feed

Recent signals

Time-series charts

🔍 Key Highlights

Clean Architecture + Domain-Driven Design

Multi-role modular services (cmd/admin, cmd/trader, cmd/customer)

Secure JWT authentication

Validated wallet system

Real-time signals + subscriptions

Production-ready directory structure

Docker-ready deployment

👥 Role Overview
👤 Customer

Signup/Login

Browse traders

Subscribe/unsubscribe

Upload KYC

View subscribed signals

Wallet operations (deposit/withdraw/history)

👨‍💼 Trader

Create & publish signals

Push live trades

Create/manage subscription plans

View subscribers

Manage profile

🛡️ Admin

Manage traders & customers

Manage subscription plans

Dashboard & analytics

System configuration

Audit logs

🧱 Architecture Overview
🧩 Clean Architecture Layers

Handlers — HTTP, validation, routing

Services / Use-Cases — Core business logic

Repositories — Data persistence

Domain Models — Pure business rules, no external dependencies

📁 Project Structure
tradeverse/
├── cmd/
│   ├── admin/
│   ├── trader/
│   └── customer/
│
├── internal/
│   ├── admin/
│   ├── trader/
│   ├── customer/
│   └── migrations/
│
├── pkg/
│   ├── auth/
│   ├── models/
│   ├── seeder/
│   └── utils/
│
├── config/
├── static/
├── templates/
└── README.md

🔄 Request Flow
Client → Router → Handler → Service → Repository → PostgreSQL

⚙️ Core Modules
🔐 Authentication

JWT-based

RBAC middleware

Token claims for role + expiry + subscription

💳 Wallet System

Deposit/Withdraw

Transaction history

Race-condition safe

Per-role actions

Ledger accuracy guaranteed

📡 Trader Module

CRUD signals

Publish live trades

Subscription plans

Subscriber management

🧾 Customer Module

Explore traders

Subscribe/unsubscribe

See signals of subscribed traders

KYC upload

Wallet summary

🛠️ Admin Module

User, trader, customer management

Subscription plans

Dashboard & analytics

📘 API Endpoints (High-Level)
Trader
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

Customer
/signup
/login
/ListTraders
/GetTraderDetails
/SubscribeToTrader
/GetSignalsFromSubscribedTraders
/kycDocument
/GetWalletSummary

Admin
/ListAdminSubscriptionPlans
/SubscribeToAdminPlan
/CancelAdminSubscription

🧰 Tech Stack

Go (Golang) — Gin framework

PostgreSQL

GORM ORM

Cron Jobs

Server-side rendered Admin UI

JWT Authentication

Clean Architecture + DDD

Docker-ready

🔧 How Internals Work
1️⃣ Market Price Fetcher

Scheduler triggers every X seconds

Calls market APIs

Normalizes & stores price data

Pushes updates to UI or cache

2️⃣ Subscription Watcher

Runs every few minutes

Deactivates expired subscriptions

Sends events/notifications

3️⃣ RBAC Engine

JWT claim inspection

Middleware checks before handler execution

4️⃣ Admin Panel

Configurable system settings

Commission & pricing

Complete user lifecycle

5️⃣ Signal Cards

Live current price

Entry/SL/Target UI formatting

Status-based color coding

6️⃣ Dashboard

Charts for:

Revenue

Subscription growth

Active signals

Market data

▶️ Running the Project
Admin Service
go run cmd/admin/main.go

Trader Service
go run cmd/trader/main.go

Customer Service
go run cmd/customer/main.go

Migrations
go run internal/migrations/main.go

Seeder
go run pkg/seeder/main.go

🔐 Security Considerations

JWT expiry & rotation

API keys managed externally (Vault/AWS Secrets Manager)

Rate limiting for market APIs

SQL injection protection

HTTPS + Nginx reverse proxy

Strong CORS policy

🚀 Deployment

Docker / Docker Compose

Kubernetes-ready

Separate worker containers (cron jobs)

Prometheus metrics

Redis for caching / pub-sub

Managed PostgreSQL

🎯 Why This Project Stands Out

Rare multi-role fintech architecture

Realistic wallet, trader signals, and subscription engine

Clean Architecture + DDD (industry standard)

Separate executables for horizontal scaling

Strong backend engineering practices

🔧 Performance

Optimized DB queries

Zero N+1 queries

Wallet consistency via transactions

Microservice-ready split

Supports future Kafka/Redis integration

🖥️ System Diagram
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

👩‍💻 Author

Fathima Sithara
Backend Developer (Golang • Microservices • Full Stack Capable)
GitHub: https://github.com/fathimasithara01