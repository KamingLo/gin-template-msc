# Go-Gin Modular Clean Architecture

A robust, scalable, and secure backend foundation built with **Golang 1.26+** and the **Gin Gonic** framework. This template implements a modular layered architecture, making it ideal for high-performance mobile-first APIs and modern web applications.

Featuring integrated **Email OTP Service**, **Google OAuth2**, **Password Recovery**, **Adaptive Rate-Limiting**, and **Enhanced CORS Security**.

## 🚀 Features

- **Modular Architecture**: Clean separation of Models, Controllers, Services, and Routes.
- **Robust Authentication**:
  - Email & Password Login.
  - OTP-based Registration (with dynamic cooldowns to prevent spam).
  - Google OAuth2 Integration.
  - Forgot & Reset Password flows.
- **Security First**:
  - JWT (JSON Web Tokens) for session management.
  - Custom CORS middleware.
  - Adaptive IP and Path-based Rate Limiting to prevent brute-force attacks.
- **Database & ORM**: PostgreSQL integration using GORM.
- **Container Ready**: Includes `docker-compose.yml` for instantly spinning up PostgreSQL
- **Developer Experience**: Pre-configured for live reloading using Air.

## 🏗️ Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go       # Application Entry Point
├── config/               # Database & Third-Party Configs (OAuth)
├── controllers/          # HTTP Request Handlers & Input Validation
├── models/               # GORM Database Schemas
├── routes/               # API Routes & Custom Middleware
├── services/             # Core Business Logic (Auth, OTP, Mail, etc.)
├── templates/            # HTML Email Templates
├── utils/                # Shared Helpers (JWT, Passwords, Standardized Responses)
├── docker-compose.yml    # Infrastructure as Code (PostgreSQL)
└── .env.example          # Environment Variable Definitions
```

## 🛠️ Getting Started

### 1. Prerequisites

- **Go 1.26+**
- **Docker & Docker Compose** (for running the database locally)
- **Air** (for live reloading): `go install github.com/air-verse/air@latest`

### 2. Project Initialization

To use this template for your own project, ensure you update the Go module name:

1. **Rename the Go Module**:
   ```bash
   go mod edit -module your-project-name
   ```

2. **Global Import Path Update**:
   Update all internal import paths from `"template/...` to `"your-project-name/...` across the project using your code editor's search and replace function.

3. **Tidy Dependencies**:
   ```bash
   go mod tidy
   ```

### 3. Environment Configuration

Copy the `.env.example` file to `.env`:
```bash
cp .env.example .env
```
Update the `.env` file with your specific database credentials, JWT secrets, Google OAuth keys, and SMTP configuration.

### 4. Infrastructure Setup

Start the required services (PostgreSQL ) using Docker Compose:
```bash
docker-compose up -d
```

### 5. Run the Application

**Development (with live reload):**
```bash
air
```

**Production (Build and Run):**
```bash
go build -o main cmd/api/main.go
./main
```

## 🔐 API Endpoints

### Authentication (`/auth`)
| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| **POST** | `/auth/otp` | Request Registration OTP | No |
| **POST** | `/auth/register` | Register new user using Email + OTP | No |
| **POST** | `/auth/login` | Authenticate using Email & Password | No |
| **GET** | `/auth/google` | Initiate Google OAuth2 flow | No |
| **GET** | `/auth/google/callback` | Google OAuth2 callback handler | No |
| **POST** | `/auth/forgot-password`| Request password reset link | No |
| **POST** | `/auth/reset-password` | Reset password using token | No |
| **GET** | `/auth/me` | Get current authenticated user details| **Yes (JWT)** |
| **GET** | `/auth/logout` | Invalidate current session | **Yes (JWT)** |

### Books (`/books`) - Example CRUD
| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| **GET** | `/books` | List all available books | No |
| **POST** | `/books/` | Create a new book entry | **Yes (JWT)** |
| **PATCH**| `/books/:id` | Update specific book details | **Yes (JWT)** |
| **DELETE**| `/books/:id` | Remove a book entry | **Yes (JWT)** |

## 🛡️ Core Middleware

- **`AuthMiddleware()`**: Validates incoming requests by verifying the attached JWT `Bearer` token.
- **`RateLimitMiddleware()`**: Protects endpoints by tracking IP and Path combinations, enforcing limits like 5 requests per 30 seconds.
- **`CORSMiddleware()`**: Enforces cross-origin policies dictated by the `ALLOWED_ORIGINS` environment variable.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
