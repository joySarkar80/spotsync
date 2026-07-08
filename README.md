Markdown# SpotSync 🚗⚡
### Smart Parking & EV Charging Reservation System

SpotSync is a robust, high-performance backend platform designed for busy hubs like airports and shopping malls. It manages parking zones and efficiently handles high-demand reservations for limited EV charging spots using enterprise-grade concurrency controls.

The system is built with **Go (Golang)**, **Echo framework**, and **GORM** (PostgreSQL), strictly enforcing Clean Architecture through a **Domain-Driven Design (DDD)** folder structure to separate concerns and guarantee atomicity.

---

## 🏗️ Architecture & Core Components

This project is organized into self-contained domain boundaries inside the `internal/domain` directory. Each domain encapsulates its own lifecycle, ensuring high maintainability and preventing database structures from leaking to the presentation layer.

### Directory Breakdown
```text
.
├── cmd/
│   └── main.go                  # App Entrypoint & Manual Dependency Injection (Wiring)
├── internal/
│   ├── auth/                    # JWT Token signing, verification, and claims processing
│   ├── config/                  # Configuration parsing & GORM PostgreSQL setup
│   ├── domain/                  # Isolated Core Business Domains
│   │   ├── user/                # User Context (Registration, Login, RBAC structures)
│   │   ├── zone/                # Parking Zones (Creation, Dynamic Spot Calculations)
│   │   └── reservation/         # Bookings (Atomic transactions, Concurrency locks)
│   ├── httpresponse/            # Uniform API JSON wrappers for successes and errors
│   ├── middlewares/             # Security layer for JWT parsing & Role verification
│   └── server/                  # Echo HTTP server initialization and route routing
└── README.md
🚨 Concurrency Critical: Solving the "EV Spot Bottleneck"When multiple drivers attempt to reserve the final remaining EV spot at the exact same millisecond, standard database queries create a race condition.SpotSync prevents overselling by wrapping capacity checks and booking creation inside a GORM Database Transaction, combined with an explicit PostgreSQL Row-Level Lock (FOR UPDATE) on the targeted parking zone record. This guarantees that only one request acquires the lock, reads the valid capacity, and completes atomically.🛠️ Technology StackLanguage: Go (Golang) v1.22 or higherWeb Framework: Echo v4 (github.com/labstack/echo/v4)ORM: GORM with PostgreSQL driver (gorm.io/gorm)Database: PostgreSQL (NeonDB, Supabase, or Aiven)Validation: Struct Validator v10 (github.com/go-playground/validator/v10)Security: JWT (github.com/golang-jwt/jwt/v5) & Bcrypt (golang.org/x/crypto/bcrypt with cost 10-12)🚀 Getting Started1. PrerequisitesEnsure you have Go 1.22+ installed and access to a working PostgreSQL cluster.2. Environment ConfigurationCreate a .env file in the root of the project:Code snippetPORT=8080
DB_URL=postgres://your_user:your_password@your_host:5432/spotsync?sslmode=require
JWT_SECRET=your_super_secure_jwt_secret_key
3. ExecutionInitialize modules, fetch external dependencies, and boot the application:Bash# Download dependencies
go mod download

# Start the server
go run cmd/main.go
The application will spin up and listen on your designated environment port (e.g., http://localhost:8080).🔐 Roles & Access Control MatrixSystem ActionDriverAdminAccount Registration & Authentication✅✅View Zones & Dynamic available_spots✅✅Reserve Parking / EV Spot✅✅View / Cancel Own Bookings✅✅Create, Update, Delete Parking Zones❌✅Modify Zone Hourly Rates❌✅View System-Wide Reservations❌✅🌐 API SpecificationsAll communications interact with standard prefixes and return uniform response structures.🔹 Authentication ModulePOST /api/v1/auth/register (Public) — Registration for accounts (driver or admin).POST /api/v1/auth/login (Public) — Verification that responds with a signed JWT token carrying claims (id, role).🔹 Parking Zones ModulePOST /api/v1/zones (Admin Only) — Provisions a new parking or EV zone.GET /api/v1/zones (Public) — Lists all parking zones featuring real-time subquery calculation for available_spots.GET /api/v1/zones/:id (Public) — Obtains details for a specific zone including dynamic availability.🔹 Reservations ModulePOST /api/v1/reservations (Authenticated) — Places a validated booking under atomic concurrency locks.GET /api/v1/reservations/my-reservations (Authenticated) — Fetches active and previous bookings specific to the requesting token user.DELETE /api/v1/reservations/:id (Authenticated) — Cancels an active reservation (rejects with 403 Forbidden if attempting to delete another user's spot).GET /api/v1/reservations (Admin Only) — Returns all system reservations using eager GORM preloads for Users and Zones.📋 Centralized Response StandardStandard Success Response (200 OK / 201 Created)JSON{
  "success": true,
  "message": "Resource processed successfully",
  "data": { ... }
}
Standard Error Response (400 / 401 / 403 / 404 / 409 / 500)JSON{
  "success": false,
  "message": "Error classification description",
  "errors": "Detailed validation rules failing or logical conflict reason"
}