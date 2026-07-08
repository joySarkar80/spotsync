# spotsyncAmm — Mango Shop REST API

A Go REST API for a mango shop. Supports user auth (JWT), mango inventory management, and order placement. Uses Echo v5, GORM, and PostgreSQL (Neon).

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Step 1 — Install Go](#step-1--install-go)
- [Step 2 — Set Up PostgreSQL (Neon)](#step-2--set-up-postgresql-neon)
- [Step 3 — Clone / Create Project](#step-3--clone--create-project)
- [Step 4 — Initialize Go Module](#step-4--initialize-go-module)
- [Step 5 — Install Dependencies](#step-5--install-dependencies)
- [Step 6 — Configure Environment](#step-6--configure-environment)
- [Step 7 — Install Air (Hot Reload)](#step-7--install-air-hot-reload)
- [Step 8 — Run the Project](#step-8--run-the-project)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [How Each Package Works](#how-each-package-works)

---

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | Language runtime |
| Git | Any | Version control |
| PostgreSQL | Any (or Neon cloud) | Database |
| Air | Latest | Hot reload (optional) |

---

## Step 1 — Install Go

1. Download from [https://go.dev/dl/](https://go.dev/dl/)
2. Run the installer
3. Verify:

```bash
go version
# go version go1.25.0 windows/amd64
```

4. Make sure `GOPATH` and `GOBIN` are in your PATH. Usually auto-set by the installer.

```bash
go env GOPATH
go env GOBIN
```

---

## Step 2 — Set Up PostgreSQL (Neon)

This project uses **Neon** — a free serverless PostgreSQL. You can also use a local PostgreSQL instance.

### Option A: Neon (recommended, free)

1. Go to [https://neon.tech](https://neon.tech) and sign up
2. Create a new project (e.g. `mangoshop`)
3. Go to **Dashboard → Connection Details**
4. Copy the connection string. It looks like:
   ```
   postgresql://user:password@ep-xxx.region.aws.neon.tech/dbname?sslmode=require
   ```
5. Paste it into your `.env` file as `DSN=...` (see Step 6)

### Option B: Local PostgreSQL

1. Install PostgreSQL from [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
2. Create a database:

```bash
psql -U postgres
CREATE DATABASE mangoshop;
\q
```

3. Your DSN will be:
   ```
   DSN="postgresql://postgres:yourpassword@localhost:5432/mangoshop?sslmode=disable"
   ```

> **Note:** GORM auto-migrates all tables on startup. You do **not** need to create tables manually.

---

## Step 3 — Clone / Create Project

If cloning an existing repo:

```bash
git clone <repo-url>
cd spotsync
```

If starting from scratch:

```bash
mkdir spotsync
cd spotsync
```

---

## Step 4 — Initialize Go Module

```bash
go mod init spotsync
```

This creates `go.mod`. The module name `spotsync` is used as the import prefix across all internal packages (e.g. `spotsync/internal/config`).

---

## Step 5 — Install Dependencies

Run these commands one by one. Each installs a specific package and adds it to `go.mod` and `go.sum`.

### HTTP Framework — Echo v5

```bash
go get github.com/labstack/echo/v5
```

Echo is the web framework. It handles routing, middleware, and JSON binding.

- Docs: [https://echo.labstack.com/](https://echo.labstack.com/)

### ORM — GORM

```bash
go get gorm.io/gorm
```

GORM is the ORM used to define models and interact with the database. It auto-migrates structs to DB tables.

- Docs: [https://gorm.io/docs/](https://gorm.io/docs/)

### GORM PostgreSQL Driver

```bash
go get gorm.io/driver/postgres
```

This is the Postgres adapter for GORM. Always install this alongside `gorm.io/gorm` when using PostgreSQL.

### JWT — golang-jwt

```bash
go get github.com/golang-jwt/jwt/v5
```

Used to generate and validate access + refresh tokens.

- Docs: [https://pkg.go.dev/github.com/golang-jwt/jwt/v5](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)

### Validator — go-playground/validator

```bash
go get github.com/go-playground/validator/v10
```

Used to validate request DTOs using struct tags like `validate:"required,email"`.

- Docs: [https://pkg.go.dev/github.com/go-playground/validator/v10](https://pkg.go.dev/github.com/go-playground/validator/v10)

### Env Loader — godotenv

```bash
go get github.com/joho/godotenv
```

Loads `.env` file into `os.Getenv()` at runtime.

- Docs: [https://pkg.go.dev/github.com/joho/godotenv](https://pkg.go.dev/github.com/joho/godotenv)

### Password Hashing — bcrypt

```bash
go get golang.org/x/crypto
```

Used to hash and compare user passwords with bcrypt.

### UUID Generator

```bash
go get github.com/google/uuid
```

Used to generate unique order codes.

### Download All at Once (if go.mod already has them)

```bash
go mod tidy
```

This downloads all missing packages and removes unused ones.

---

## Step 6 — Configure Environment

Create a `.env` file in the project root:

```bash
# .env
DSN="postgresql://user:password@host/dbname?sslmode=require"
PORT=8080
JWT_SECRET=your_secret_key_here
```

| Variable | Description |
|----------|-------------|
| `DSN` | Full PostgreSQL connection string |
| `PORT` | Port the server listens on |
| `JWT_SECRET` | Secret used to sign JWT tokens — change this in production |

> **Warning:** Never commit `.env` to git. It is already in `.gitignore`.

---

## Step 7 — Install Air (Hot Reload)

Air watches your `.go` files and restarts the server on changes. Optional but recommended for development.

```bash
go install github.com/air-verse/air@latest
```

Verify:

```bash
air -v
```

If `air` is not found, add Go's bin directory to PATH:

- **Windows:** Add `%USERPROFILE%\go\bin` to your system PATH
- **Linux/Mac:** Add `$HOME/go/bin` to your `~/.bashrc` or `~/.zshrc`

The `.air.toml` config file is already in the project. It builds to `./tmp/main.exe` and watches `.go`, `.html`, `.tmpl` files.

---

## Step 8 — Run the Project

### With Air (hot reload)

```bash
air
```

### Without Air (standard run)

```bash
go run ./cmd/main.go
```

### Build binary

```bash
go build -o ./tmp/main.exe ./cmd/main.go
./tmp/main.exe
```

### Expected output

```
Database connection successful
⇨ http server started on [::]:8080
```

### Verify server is running

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Project Structure

```
spotsync/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── config/
│   │   ├── config.go            # Loads .env into Config struct
│   │   └── db.go                # Opens GORM DB connection
│   ├── auth/
│   │   └── jwt.go               # JWT service (generate + validate tokens)
│   ├── middlewares/
│   │   └── auth.go              # JWT auth middleware for protected routes
│   ├── httpresponse/
│   │   └── error.go             # Standard error response struct
│   ├── server/
│   │   └── http.go              # Creates Echo, auto-migrates DB, registers routes
│   └── domain/
│       ├── user/                # User registration, login, refresh, /me
│       │   ├── entity.go        # User model + bcrypt methods
│       │   ├── repository.go    # DB queries
│       │   ├── service.go       # Business logic
│       │   ├── handler.go       # HTTP handlers
│       │   ├── register.go      # Route registration
│       │   └── dto/             # Request/response structs
│       ├── mango/               # Mango inventory
│       │   ├── entity.go        # Mango model
│       │   ├── repository.go
│       │   ├── service.go
│       │   ├── handler.go
│       │   ├── register.go
│       │   └── dto/
│       └── order/               # Order placement
│           ├── entity.go        # Order model (pending/confirmed/cancelled)
│           ├── repository.go
│           ├── service.go
│           ├── handler.go
│           ├── register.go
│           └── dto/
├── .env                         # Environment variables (not in git)
├── .air.toml                    # Air hot reload config
├── go.mod                       # Module definition + dependency list
├── go.sum                       # Dependency checksums
└── mangoshop.postman_collection.json  # Postman collection for testing
```

---

## API Endpoints

### Health

| Method | URL | Auth |
|--------|-----|------|
| GET | `/health` | No |

### Auth (`/api/v1/auth`)

| Method | URL | Auth | Description |
|--------|-----|------|-------------|
| POST | `/api/v1/auth/register` | No | Create account |
| POST | `/api/v1/auth/login` | No | Login, returns access + refresh token |
| POST | `/api/v1/auth/refresh` | No | Get new access token using refresh token |
| GET | `/api/v1/auth/me` | Yes | Get current user info |

### Mangoes (`/api/v1/mangoes`)

| Method | URL | Auth | Description |
|--------|-----|------|-------------|
| GET | `/api/v1/mangoes` | No | List all mangoes |
| GET | `/api/v1/mangoes/:id` | No | Get mango by ID |
| POST | `/api/v1/mangoes` | Yes | Create new mango |
| PATCH | `/api/v1/mangoes/:id` | Yes | Update mango |

### Orders (`/api/v1/orders`)

| Method | URL | Auth | Description |
|--------|-----|------|-------------|
| POST | `/api/v1/orders` | Yes | Place an order |
| GET | `/api/v1/orders/me` | Yes | Get my orders |

**Auth** = requires `Authorization: Bearer <access_token>` header.

---

## How Each Package Works

### `internal/config`

Reads `.env` with `godotenv` and returns a `Config` struct. `db.go` opens a GORM connection using the DSN. Called once at startup in `main.go`.

### `internal/auth`

`JWTService` interface with two implementations: `GenerateAccessToken` (15 min TTL) and `GenerateRefreshToken` (7 days TTL). Tokens are signed with HS256.

### `internal/middlewares`

`AuthMiddleware` validates the Bearer token from the `Authorization` header and sets `user_id`, `user_email`, `user_name` in Echo's context for downstream handlers.

### `internal/server`

Creates the Echo instance, attaches the custom validator (wraps `go-playground/validator`), calls `db.AutoMigrate` (creates tables if they don't exist), then calls each domain's `RegisterRoutes`.

### `internal/domain/user`

- Register: hashes password with bcrypt, stores user
- Login: compares password, returns access + refresh tokens
- Refresh: validates refresh token, issues new access token
- Me: reads user info from context (set by middleware)

### `internal/domain/mango`

CRUD for mango inventory. Create and Update require auth. Get routes are public.

### `internal/domain/order`

Creates an order by looking up the mango, checking stock, calculating total price, decrementing stock, and generating a unique order code (UUID). All routes require auth.

---

## Testing with Postman

Import `mangoshop.postman_collection.json` into Postman:

1. Open Postman → Import → select the file
2. Register a user first (`POST /api/v1/auth/register`)
3. Login to get tokens (`POST /api/v1/auth/login`)
4. Set the access token in the `Authorization` header as `Bearer <token>` for protected routes

---

## Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `Error loading .env file` | `.env` file missing | Create `.env` in project root |
| `failed to connect database` | Wrong DSN | Check DSN in `.env` |
| `air: command not found` | Go bin not in PATH | Add `~/go/bin` to PATH |
| `port already in use` | Another process on port 8080 | Change `PORT` in `.env` or kill the process |
