# Go-Gin Modular Clean Architecture (v2.0)

A robust, scalable, and secure backend foundation built with **Golang 1.26+** and the **Gin Gonic** framework. This template implements a modular layered architecture, making it ideal for high-performance mobile-first APIs and Progressive Web Apps (PWA).

New in this version:  **Integrated OTP Service** ,  **Google OAuth2** ,  **Adaptive Rate-Limiting** , and  **Enhanced CORS Security** .

## 🏗️ Architecture Overview

The project follows a strict separation of concerns to ensure the codebase remains maintainable as it grows:

* **Models** : GORM-based data structures and database schemas.
* **Services** : The "Brain" layer. Handles business logic, OTP cooldown calculations, and third-party integrations (Email/OAuth).
* **Controllers** : HTTP entry points. Responsible for input validation and returning standardized JSON.
* **Routes & Middleware** : Modular routing with built-in security (JWT, Rate Limiting, CORS).
* **Utils/Config** : Shared helpers for security, environment loading, and database connectivity.

---

## 📂 Project Structure

**Plaintext**

```
template/
├── cmd/api/main.go        # Entry point
├── config/                # DB, Google OAuth
├── controllers/          # HTTP Handlers (using utils for responses)
├── models/               # Database Schemas
├── routes/               # Modular routes & Custom Middleware
├── services/             # Business Logic (OTP, Mail, Auth, Books)
├── templates/            # HTML Email Templates (OTP)
├── utils/                # Standardized Response & Security Helpers
└── .env                  # Environment Variables
```

---

## 🛠️ Getting Started

### 1. Prerequisites

* **Go 1.26+** installed.
* **PostgreSQL** instance.
* **Air** for live reloading: `go install github.com/air-verse/air@latest`

### 2. Project Initialization (CRITICAL ⚠️)

To ensure internal imports work correctly, you must perform these three steps:

1. **Rename the Go Module** :
   **Bash**

```
   go mod edit -module your-project-name
```

2. **Global Refactor (Search & Replace)** :
   Since all internal imports use the `template/` prefix, you must replace them.

* **In VS Code** : Press `Ctrl + Shift + H`.
* **Search** : `"template/`
* **Replace** : `"your-project-name/`
* *Note: Including the opening quote ensures you only modify import paths, not your logic.*

3. **Tidy Dependencies** :
   **Bash**

```
   go mod tidy
```

---

### 3. Configuration (.env)

Create a `.env` file in the root directory. This template is pre-configured for **Gmail SMTP** and  **Google OAuth** :

**Cuplikan kode**

```
# Server Config
GIN_MODE=debug
MACHINE_ID=dc1
PORT=8000

# Security & CORS
COOKIE_DOMAIN=localhost
ALLOWED_ORIGINS=http://127.0.0.1:3000,http://localhost:3000

# Database (PostgreSQL)
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_db_name
DB_PORT=5432
DB_SSLMODE=require

# JWT Configuration
JWT_SECRET=your-32-or-64-character-secret
JWT_EXPIRES=enable
JWT_EXPIRES_IN=24

# Google OAuth2
GOOGLE_CLIENT_ID=your_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_secret
GOOGLE_CALLBACK_URL=http://localhost:8000/auth/google/callback

# Frontend Redirects
OAUTH_FRONTEND_URL=http://127.0.0.1:3000/auth/login
SUCCESS_FRONTEND_URL=http://127.0.0.1:3000/auth/callback
SESSION_SECRET=your-random-session-secret

# SMTP / Mail Settings
MAIL_MAILER=smtp
MAIL_HOST=smtp.gmail.com
MAIL_PORT=465
MAIL_USERNAME=your_email@gmail.com
MAIL_PASSWORD=your_google_app_password
MAIL_ENCRYPTION=ssl
MAIL_FROM_ADDRESS=noreply@yourdomain.com
MAIL_FROM_NAME="App Verifier"
```

---

## 🔐 API Endpoints

| **Method** | **Endpoint** | **Description**      | **Auth Required** |
| ---------------- | ------------------ | -------------------------- | ----------------------- |
| **POST**   | `/auth/otp`      | Request Registration OTP   | No                      |
| **POST**   | `/auth/register` | Register using Email + OTP | No                      |
| **POST**   | `/auth/login`    | Email & Password Login     | No                      |
| **GET**    | `/auth/google`   | Trigger Google OAuth       | No                      |
| **GET**    | `/books`         | List all books             | No                      |
| **GET**    | `/books/:id`     | Get specific book details  | No                      |
| **POST**   | `/books`         | Add a new book             | **Yes (JWT)**     |
| **PATCH**  | `/books/:id`     | Partial update book        | **Yes (JWT)**     |
| **DELETE** | `/books/:id`     | Remove a book              | **Yes (JWT)**     |

---

## 🛡️ Security & Features

* **Adaptive Rate Limiting** : Protects your API from brute-force. Hits trigger a 30-second lockout with a dynamic "retry-in" message.
* **OTP Cooldown Logic** : Built-in anti-spam for SMTP. Wait times increase (30s → 1m → 5m → 1h) based on request frequency.
* **Clean Error Handling** : Uses `utils.SendError` and `utils.SendSuccess` to ensure mobile clients always receive a consistent JSON structure.
* **Background Workers** : Email sending is handled in a separate goroutine with panic recovery, ensuring your API remains lightning-fast.
* **Mobile-First Design** : The OTP email template is fully responsive and uses system fonts to look native on iOS and Android.

## 🚀 Run the App

**Development (with live reload):**

**Bash**

```
air
```

**Production (Binary build):**

**Bash**

```
go build -o main cmd/api/main.go
./main
```
