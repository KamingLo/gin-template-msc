# 🚀 Go-Gin Modular Clean Architecture

A robust and scalable production-ready backend template built with **Golang** and the **Gin Gonic** framework. This project implements a modular layered architecture to ensure maintainability and ease of testing.

## 🏗️ Architecture Overview

This project follows the **Models-Services-Controllers** pattern, ensuring that each component has a single responsibility:

* **Models**: Defines the data structures and GORM database schemas.
* **Services**: The "Brain" of the application. Handles business logic, complex validations, and data processing.
* **Controllers**: Entry point for HTTP requests. Handles input binding and returns appropriate JSON responses.
* **Routes**: Modular routing system split by domain (Auth, Books, etc.) to keep the codebase clean.
* **Utils/Config**: Contains shared helpers (JWT, Hashing) and database connection logic.

---

## 📂 Project Structure

```text
belajar-go/
├── cmd/api/main.go       # Application Entry Point
├── config/               # Database & Environment Configuration
├── controllers/          # HTTP Request Handlers
├── models/               # Data Structures & GORM Schemas
├── routes/               # Modular Routing & Middleware
├── services/             # Business Logic Layer
├── utils/                # Security Helpers (JWT, Bcrypt)
├── .env                  # Environment Variables (Secrets)
└── .gitignore            # Files to be excluded from Git

```

---

## 🛠️ Getting Started (Development Mode)

### 1. Prerequisites

* **Go 1.21+** installed.
* A running **PostgreSQL** instance (Local or Supabase).
* **Air** for live reloading (Recommended):

```bash
go install github.com/air-verse/air@latest

```

### 2. Configuration

Create a `.env` file in the root directory:

```env
GIN_MODE=debug
PORT=4000

DB_HOST=your-database-host
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_PORT=5432
DB_SSLMODE=require

JWT_SECRET=your-secure-random-string
JWT_EXPIRES=enable
JWT_EXPIRES_IN=2

```

### 3. Run the App

Untuk menjalankan server dengan fitur **live reload** (setiap ada perubahan kode, server otomatis restart):

```bash
air

```

*Atau secara manual:*

```bash
go run cmd/api/main.go

```

---

## 🚀 Deployment (Binary)

Untuk mendeploy ke VPS atau PaaS (Railway, Render, dll) tanpa container:

1. **Build binary**:

```bash
go build -o main cmd/api/main.go

```

2. **Set Environment Variables**: Pastikan variabel di `.env` sudah diatur pada server/dashboard provider.
3. **Jalankan**: `./main`

---

## 🔐 API Endpoints

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| **POST** | `/auth/register` | Register a new user | No |
| **POST** | `/auth/login` | Login & receive JWT | No |
| **GET** | `/books` | List all books | No |
| **POST** | `/books` | Add a new book | Yes (JWT) |
| **PATCH** | `/books/:id` | Partial update book | Yes (JWT) |
| **DELETE** | `/books/:id` | Remove a book | Yes (JWT) |

---

## 🛡️ Security Features

* **Password Hashing**: Menggunakan `bcrypt` untuk menyimpan password secara aman.
* **JWT Authentication**: Autentikasi stateless menggunakan JSON Web Tokens.
* **Environment Safety**: Kredensial kritikal tidak pernah di-hardcode dan dikecualikan dari Git via `.gitignore`.