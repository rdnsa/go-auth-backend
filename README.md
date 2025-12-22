# Go Auth Backend – Clean Architecture + MongoDB + JWT

A simple yet production-ready authentication API built with Go, implementing Clean Architecture, JWT authentication, and MongoDB as the primary database.

## ✨ Features

User registration

Login with JWT token (configurable expiration)

Passwords automatically hashed using bcrypt

Logging for successful and failed register/login events

Input validation

.env configuration support

Clean Architecture structure (entity → usecase → repository → handler)

## 🧩 Tech Stack

Go 1.23+

Gin Web Framework

MongoDB (mongo-driver)

JWT (HS256)

bcrypt for password hashing

godotenv for environment management

📁 Project Structure
go-auth-backend/
├── cmd/api/main.go                 # Entry point
├── internal/
│   ├── config/                     # Config & .env setup
│   ├── dto/                        # Request/response DTOs
│   ├── entity/                     # Domain models
│   ├── handler/http/               # Gin HTTP handlers
│   ├── repository/mongodb/         # Database layer
│   └── usecase/                    # Business logic
├── pkg/
│   ├── jwt/                        # JWT helper
│   └── password/                   # bcrypt helper
├── .env
├── .env.example
├── go.mod
└── README.md

⚙️ Run Locally
1. Start MongoDB

Easiest way using Docker:

docker run -d -p 27017:27017 --name mongo-auth mongo:latest

2. Setup Project
# Clone repository & enter directory
git clone https://github.com/username/go-auth-backend.git
cd go-auth-backend

# Copy environment file
cp .env.example .env

# Install dependencies
go mod tidy

# Run server
go run cmd/api/main.go


Server runs at:
👉 http://localhost:8080

🔑 Test API
Register
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi",
    "email": "budi@gmail.com",
    "password": "rahasia123"
  }'

Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@gmail.com",
    "password": "rahasia123"
  }'


Response Example:

{
  "message": "Login successful!",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxx",
    "user": {
      "id": "6759a1b2c3d4e5f6a7b8c9d0",
      "name": "Budi",
      "email": "budi@gmail.com"
    }
  }
}

## 🧠 View Data in MongoDB Compass

Open MongoDB Compass

Connect to:

mongodb://localhost:27017


Select database: auth_db

Open collection: users

## 🧾 Example .env
APP_PORT=8080
JWT_SECRET=super-secret-jwt-key-1234567890abcdef
JWT_EXPIRED_HOURS=72

MONGO_URI=mongodb://localhost:27017
MONGO_DB=auth_db
MONGO_USER_COLLECTION=users

## 🧭 Upcoming Features

Protected routes + JWT middleware (/profile)

Refresh token support

Docker Compose setup

Unit & integration tests

Free deployment options (Railway, Render, Fly.io)

### 📜 License

MIT © 2025
Built with Go and a spirit of learning.
