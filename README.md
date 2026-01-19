# SPSC loanEasy v1.0

ระบบสินเชื่อสหกรณ์ SPSC - Backend API

## 🚀 Tech Stack

- **Framework:** Go Fiber v2
- **Database:** MySQL + GORM
- **Documentation:** Swagger
- **Architecture:** Clean Architecture

## 📁 Project Structure

```
spsc-loaneasy/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point
│
├── internal/
│   ├── adapters/
│   │   ├── http/
│   │   │   ├── handlers/              # HTTP handlers
│   │   │   ├── middleware/            # HTTP middlewares
│   │   │   └── routes/                # Route definitions
│   │   │
│   │   └── persistence/
│   │       ├── models/                # GORM models
│   │       └── repositories/          # Repository implementations
│   │
│   ├── config/
│   │   ├── config.go                  # Configuration loader
│   │   ├── database.go                # Database connection
│   │   └── seeder.go                  # Database seeder
│   │
│   └── core/
│       ├── domain/                    # Domain entities & errors
│       └── services/                  # Service interfaces
│
├── docs/                              # Swagger documentation
├── .env                               # Environment variables
├── .env.example                       # Environment template
├── Makefile                           # Build commands
├── go.mod
├── GUIDE.md                           # Installation guide
└── README.md
```

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Request                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapters (HTTP Layer)                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Handlers   │  │ Middleware  │  │      Routes         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Core (Business Logic)                     │
│  ┌─────────────────────┐  ┌─────────────────────────────┐   │
│  │       Domain        │  │         Services            │   │
│  │  (Entities/Errors)  │  │   (Business Logic)          │   │
│  └─────────────────────┘  └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Adapters (Persistence Layer)                 │
│  ┌─────────────────────┐  ┌─────────────────────────────┐   │
│  │       Models        │  │       Repositories          │   │
│  │   (GORM Models)     │  │    (Data Access)            │   │
│  └─────────────────────┘  └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        Database                             │
│  ┌─────────────────────┐  ┌─────────────────────────────┐   │
│  │   flommast (R/O)    │  │   users, refresh_tokens     │   │
│  │     (Legacy)        │  │        (New)                │   │
│  └─────────────────────┘  └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## ⚙️ Configuration

ใช้ไฟล์ `.env` ไฟล์เดียวสำหรับทั้ง Dev และ Prod:

```env
APP_MODE=dev  # dev | prod

# DEV config uses DEV_ prefix
DEV_DB_HOST=localhost

# PROD config uses PROD_ prefix  
PROD_DB_HOST=production-server
```

## 🛠️ Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Setup database
mysql -u root -p -e "CREATE DATABASE spsccoop_webt2app"

# 3. Configure .env
cp .env.example .env
# Edit .env with your settings

# 4. Run server
make run-dev
```

## 📚 API Documentation

Swagger UI: http://localhost:3000/swagger/index.html

## 📝 License

Apache 2.0
