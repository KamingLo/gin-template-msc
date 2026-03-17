# 🚀 Go-Gin Modular Clean Architecture

A robust and scalable production-ready backend template built with **Golang** and the **Gin Gonic** framework. This project implements a modular layered architecture to ensure maintainability and ease of testing.

New Feature : Adding rate-limiter and cors for security
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
template/
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

* **Go 1.26+** installed.
* A running **PostgreSQL** instance.
* **Air** for live reloading (Recommended):
```bash
go install github.com/air-verse/air@latest

```



### 2. Project Initialization (IMPORTANT ⚠️)

After cloning the repository, you **must** rename the module to match your project name. Failure to do so will result in broken internal import paths.

1. **Rename the Go Module**:
Run this command in your terminal to update `go.mod`:
```bash
go mod edit -module your-project-name

```


2. **Global Refactor (Search & Replace)**:
Since all internal imports use the `template/` prefix, you need to replace them globally:
* **VS Code**: Press `Ctrl + Shift + H`.
* **Search**: `template`
* **Replace**: `your-project-name`
* *Note: Include the double quotes to ensure only import statements are modified.*


3. **Tidy Dependencies**:
```bash
go mod tidy

```



### 3. Configuration

Create a `.env` file in the root directory by copying the example:

```bash
cp .env.example .env

```

Ensure your `.env` contains the correct database credentials:

```env
GIN_MODE=debug
PORT=4000

DB_HOST=127.0.0.1
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_db_name
DB_PORT=5432
DB_SSLMODE=disable

JWT_SECRET=your-secure-random-string
JWT_EXPIRES_IN=2

```

### 4. Run the App

To run with **live reload** (automatic restart on code changes):

```bash
air

```

*Or manually:*

```bash
go run cmd/api/main.go

```

---

## 🚀 Deployment (Binary)

To deploy to a VPS or PaaS (Railway, Render, etc.) without containers:

1. **Build binary**:
```bash
go build -o main cmd/api/main.go

```


2. **Set Environment Variables**: Ensure the variables in `.env` are set in your server/dashboard provider.
3. **Run**: `./main`

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

* **Password Hashing**: Uses `bcrypt` for secure storage.
* **JWT Authentication**: Stateless authentication using JSON Web Tokens.
* **Environment Safety**: Critical credentials are never hardcoded and are excluded from Git via `.gitignore`.