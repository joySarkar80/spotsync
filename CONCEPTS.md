# Project Concepts & Keywords

A plain-English guide to every keyword, folder, file, and pattern in this project. Read this once and the entire codebase will make sense.

---

## Table of Contents

- [File & Folder Names](#file--folder-names)
  - [cmd/](#cmd)
  - [internal/](#internal)
  - [go.mod](#gomod)
  - [go.sum](#gosum)
  - [.air.toml](#airtoml)
  - [.env](#env)
  - [.gitignore](#gitignore)
- [Go Language Keywords](#go-language-keywords)
  - [package](#package)
  - [struct](#struct)
  - [interface](#interface)
  - [func](#func)
  - [Pointer vs Value receiver](#pointer-vs-value-receiver)
  - [Struct Tags](#struct-tags)
- [Project Architecture Patterns](#project-architecture-patterns)
  - [Domain](#domain)
  - [Entity](#entity)
  - [DTO (Data Transfer Object)](#dto-data-transfer-object)
  - [Repository](#repository)
  - [Service](#service)
  - [Handler](#handler)
  - [Middleware](#middleware)
  - [Register (Route Registration)](#register-route-registration)
- [GORM Concepts](#gorm-concepts)
  - [gorm.Model](#gormmodel)
  - [AutoMigrate](#automigrate)
  - [Struct Tags for GORM](#struct-tags-for-gorm)
- [Request Flow (How a request travels)](#request-flow)
- [Auth Flow (JWT)](#auth-flow-jwt)

---

## File & Folder Names

### `cmd/`

**What:** Standard Go convention for the application entry point(s).

**Why:** Go projects can have multiple runnable programs. Each one lives in `cmd/<name>/main.go`. Our project has one: `cmd/main.go`.

```
cmd/
└── main.go   ← only file. starts the whole app.
```

`main.go` does exactly 3 things:
1. Load config from `.env`
2. Connect to database
3. Start the HTTP server

```go
func main() {
    cfg := config.LoadEnv()        // step 1
    db  := config.ConnectDatabase(cfg)  // step 2
    server.Start(db, cfg)          // step 3
}
```

Everything else is wired inside `server.Start`.

---

### `internal/`

**What:** A special Go directory name. Code inside `internal/` **cannot be imported by any outside project**.

**Why:** Enforces encapsulation. If this project were a library, consumers couldn't accidentally use our internal plumbing. Only code inside `haddibanga/` can import `haddibanga/internal/...`.

```
internal/
├── config/       ← env + DB setup
├── auth/         ← JWT logic
├── middlewares/  ← HTTP middleware
├── httpresponse/ ← shared error struct
├── server/       ← Echo setup + route wiring
└── domain/       ← business logic split by feature
    ├── user/
    ├── mango/
    └── order/
```

---

### `go.mod`

**What:** The module definition file. Like `package.json` in Node.js.

**Why:** Tells Go:
- What this module is called (`module haddibanga`)
- Which Go version to use (`go 1.25.0`)
- What external packages are needed (`require` block)

```go
module haddibanga      ← this is the import prefix for all internal packages

go 1.25.0

require (
    github.com/labstack/echo/v5  v5.1.1
    gorm.io/gorm                 v1.25.11
    ...
)
```

When you run `go get github.com/some/package`, it gets added here automatically.

---

### `go.sum`

**What:** Auto-generated checksum file. Like `package-lock.json`.

**Why:** Locks the exact version + hash of every dependency so builds are reproducible and tamper-proof. **You never edit this manually.** `go mod tidy` manages it.

---

### `.air.toml`

**What:** Config file for [Air](https://github.com/air-verse/air) — a hot reload tool for Go.

**Why:** Go is a compiled language. Without Air, you'd need to manually `go build` + restart the server after every code change. Air watches your files, rebuilds, and restarts automatically.

Key settings in our `.air.toml`:

```toml
cmd = "go build -o ./tmp/main.exe ./cmd/main.go"   ← build command
bin = "./tmp/main.exe"                              ← what to run after build
include_ext = ["go", "tpl", "tmpl", "html"]        ← watch these file types
exclude_dir = ["assets", "tmp", "vendor"]           ← ignore these folders
```

Run with just: `air`

---

### `.env`

**What:** Environment variables file. Keeps secrets out of source code.

**Why:** Credentials (DB password, JWT secret) must never be hardcoded in `.go` files or committed to git. `.env` is loaded at runtime by `godotenv` and stays local-only.

```env
DSN="postgresql://user:pass@host/db"
PORT=8080
JWT_SECRET=your_secret
```

Loaded in `internal/config/config.go`:
```go
godotenv.Load()           // reads .env into OS environment
os.Getenv("PORT")         // reads individual values
```

---

### `.gitignore`

**What:** Tells git which files/folders to never track.

**Why:** Some files must never go to GitHub — secrets (`.env`), compiled binaries (`*.exe`), build output (`tmp/`), and local tool settings (`.claude/`).

Our `.gitignore`:
```
.env       ← database credentials, JWT secret
tmp/       ← compiled binaries built by Air
*.exe      ← Windows executables
go.sum     ← auto-generated, can be regenerated
.claude/   ← local AI assistant settings
```

---

## Go Language Keywords

### `package`

**What:** Every `.go` file starts with `package <name>`. Groups related files together.

**Rules:**
- All files in the same folder must have the same package name
- The special package `main` is the entry point — only `package main` can be run directly
- Everything else uses the folder name as the package name

```go
package main        // in cmd/main.go — runnable
package config      // in internal/config/ — importable library
package user        // in internal/domain/user/ — importable library
```

---

### `struct`

**What:** A custom data type that groups related fields. Like a class in other languages, but without inheritance.

**Why:** Represents real-world objects — a User, a Mango, an Order.

```go
type User struct {
    gorm.Model                // embedded: gives ID, CreatedAt, UpdatedAt, DeletedAt
    Name     string           // a text field
    Email    string           // another text field
    Password string           // another text field
}
```

Create a value from it:
```go
user := User{
    Name:  "Alice",
    Email: "alice@example.com",
}
```

---

### `interface`

**What:** Defines a set of method signatures. Any struct that has those methods automatically satisfies the interface.

**Why:** Decouples code. A service doesn't need to know *how* the repository works — only *what* it can do.

```go
// defines what a Repository must be able to do
type Repository interface {
    CreateUser(user *User) error
    GetUserByEmail(email string) (*User, error)
}

// this struct satisfies Repository because it has both methods
type repository struct {
    db *gorm.DB
}
func (r *repository) CreateUser(user *User) error { ... }
func (r *repository) GetUserByEmail(email string) (*User, error) { ... }
```

The service only depends on the interface:
```go
type service struct {
    repo Repository   // ← interface, not the concrete struct
}
```

This means you could swap the real DB for a fake one in tests without changing the service.

---

### `func`

**What:** A function. Can be standalone or attached to a struct (called a **method**).

```go
// standalone function
func Add(a, b int) int {
    return a + b
}

// method on a struct (service is the receiver)
func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {
    ...
}
```

---

### Pointer vs Value Receiver

**What:** The difference between `(s *service)` and `(s service)`.

| Syntax | Name | Meaning |
|--------|------|---------|
| `(s *service)` | pointer receiver | `s` is a reference — mutations inside the method affect the original |
| `(s service)` | value receiver | `s` is a copy — mutations don't affect the original |

**Rule of thumb:** Use pointer receivers (`*`) almost always, especially when the struct holds state (like a DB connection). This project uses pointer receivers everywhere.

---

### Struct Tags

**What:** Backtick annotations on struct fields that add metadata. Read by libraries at runtime using reflection.

```go
type CreateRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}
```

| Tag | Library | What it does |
|-----|---------|-------------|
| `json:"name"` | `encoding/json` | JSON key name when marshaling/unmarshaling |
| `json:"id,omitempty"` | `encoding/json` | Skip field in JSON output if value is zero/empty |
| `validate:"required"` | `go-playground/validator` | Field must not be empty |
| `validate:"required,email"` | `go-playground/validator` | Must not be empty AND must be valid email format |
| `validate:"min=6"` | `go-playground/validator` | String must be at least 6 characters |
| `gorm:"type:varchar(100)"` | `gorm` | SQL column type |
| `gorm:"uniqueIndex"` | `gorm` | Add UNIQUE index on this column in DB |
| `gorm:"not null"` | `gorm` | Column is NOT NULL in DB |

---

## Project Architecture Patterns

This project uses **layered architecture** inside each domain. Every domain (`user`, `mango`, `order`) has the same 5 layers:

```
entity.go      → data shape (DB model)
repository.go  → database operations
service.go     → business logic
handler.go     → HTTP input/output
register.go    → wire everything together + define routes
dto/           → request and response shapes
```

### Domain

**What:** One feature area. `user` handles auth. `mango` handles inventory. `order` handles purchases.

**Why:** Keeps code organized by feature instead of by type. All user-related code lives in `internal/domain/user/`, not scattered across the project.

---

### Entity

**What:** The database model struct. Maps directly to a DB table.

```go
// internal/domain/user/entity.go
type User struct {
    gorm.Model          // auto-adds: id, created_at, updated_at, deleted_at
    Name     string
    Email    string
    Password string
}
```

GORM turns this struct into a `users` table automatically (via `AutoMigrate`).

---

### DTO (Data Transfer Object)

**What:** Structs used to carry data *into* the API (request) or *out of* the API (response). Separate from the entity.

**Why:** You never expose the raw entity to the outside world. The `User` entity has a `Password` field — you don't want that in API responses. DTOs let you control exactly what comes in and goes out.

```
dto/
├── request.go   ← what the client sends us
└── response.go  ← what we send back to the client
```

Example:
```go
// request: what we accept from the client
type CreateRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// response: what we send back (no Password field)
type Response struct {
    ID    uint   `json:"id,omitempty"`
    Name  string `json:"name,omitempty"`
    Email string `json:"email,omitempty"`
}
```

---

### Repository

**What:** The only layer that talks to the database. Wraps raw GORM calls behind a clean interface.

**Why:** Business logic (service) should not contain SQL/GORM calls. If you switch from PostgreSQL to MySQL, you only change the repository — nothing else.

```go
// interface defines the contract
type Repository interface {
    CreateUser(user *User) error
    GetUserByEmail(email string) (*User, error)
}

// concrete implementation uses GORM
func (r *repository) CreateUser(user *User) error {
    return r.db.Create(user).Error
}
```

---

### Service

**What:** The business logic layer. Orchestrates repository calls, applies rules, and builds responses.

**Why:** Keeps handlers thin and testable. The handler doesn't decide what "login" means — the service does.

```go
func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
    user, _ := s.repo.GetUserByEmail(req.Email)  // ask repository
    if user == nil {
        return nil, ErrInvalidCredentials         // business rule: no user = invalid
    }
    user.checkPassword(req.Password)              // business rule: check password
    token, _ := s.jwtService.GenerateAccessToken(...)  // generate JWT
    return &dto.Response{AccessToken: token}, nil // build response
}
```

---

### Handler

**What:** Sits at the HTTP boundary. Reads the request, calls the service, writes the response.

**Why:** Separates HTTP concerns (parsing JSON, writing status codes) from business logic.

```go
func (h *handler) LoginUser(c *echo.Context) error {
    var req dto.LoginRequest
    c.Bind(&req)          // 1. parse incoming JSON into req struct
    c.Validate(&req)      // 2. validate fields (required, email format, etc.)
    response, err := h.service.LoginUser(req)  // 3. call service
    return c.JSON(200, response)               // 4. write JSON response
}
```

---

### Middleware

**What:** A function that runs *before* a handler. Can inspect/modify the request or block it entirely.

**Why:** Auth logic would be duplicated in every protected handler. Middleware runs once and sets user info in context.

```go
// AuthMiddleware runs before any protected route handler
func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c *echo.Context) error {
            // 1. read Authorization header
            // 2. validate token
            // 3. if invalid → return 401, stop here
            // 4. if valid → set user_id in context, continue to handler
            c.Set("user_id", claims.UserID)
            return next(c)   // ← calls the actual route handler
        }
    }
}
```

The handler then reads from context:
```go
userId := c.Get("user_id").(uint)
```

---

### Register (Route Registration)

**What:** Each domain has a `register.go` that wires up the repository → service → handler chain and attaches routes to Echo.

**Why:** Keeps `server/http.go` clean. It just calls `user.RegisterRoutes(...)`, `mango.RegisterRoutes(...)`, etc.

```go
// internal/domain/user/register.go
func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
    repo    := NewRepository(db)          // create repository
    svc     := NewService(repo, jwt)      // inject repo into service
    handler := NewHandler(svc)            // inject service into handler

    api := e.Group("/api/v1/auth")
    api.POST("/register", handler.CreateUser)   // attach handler to route
    api.POST("/login",    handler.LoginUser)
}
```

This pattern is called **dependency injection** — each layer receives its dependencies instead of creating them itself.

---

## GORM Concepts

### `gorm.Model`

**What:** A built-in GORM struct you embed in your entities. Adds 4 fields automatically.

```go
type User struct {
    gorm.Model   // ← this line adds all 4 fields below
    Name  string
    Email string
}

// gorm.Model expands to:
// ID        uint           → auto-increment primary key
// CreatedAt time.Time      → set automatically on insert
// UpdatedAt time.Time      → updated automatically on every save
// DeletedAt gorm.DeletedAt → soft delete (sets timestamp instead of actually deleting)
```

**Soft delete:** When you call `db.Delete(&user)`, GORM does NOT run `DELETE FROM users`. It sets `deleted_at = now()`. The row stays in the DB but is invisible to normal queries. This is safe — you can recover deleted data.

---

### AutoMigrate

**What:** GORM reads your struct definitions and creates/updates DB tables to match.

```go
// in internal/server/http.go
db.AutoMigrate(&user.User{}, &mango.Mango{}, &order.Order{})
```

This runs on every startup. It:
- Creates the table if it doesn't exist
- Adds new columns if you added fields to the struct
- Does **not** delete columns or data

You never write `CREATE TABLE` SQL manually in this project.

---

### Struct Tags for GORM

```go
type User struct {
    gorm.Model
    Name  string `gorm:"type:varchar(100);not null"`
    Email string `gorm:"type:varchar(255);uniqueIndex;not null"`
}
```

| Tag | DB effect |
|-----|-----------|
| `type:varchar(100)` | Column type is VARCHAR(100) |
| `not null` | Column has NOT NULL constraint |
| `uniqueIndex` | Adds a UNIQUE index — duplicate emails rejected by DB |

---

## Request Flow

How a single HTTP request travels through the entire project:

```
Client (Postman / browser)
    │
    │  POST /api/v1/auth/login
    │  Body: { "email": "...", "password": "..." }
    ▼
Echo Router  (internal/server/http.go)
    │  matches route, calls registered handler
    ▼
Middleware (if route is protected — checks JWT token)
    │  valid token → sets user_id in context
    │  invalid token → returns 401, stops here
    ▼
Handler  (internal/domain/user/handler.go)
    │  1. c.Bind(&req)     → parse JSON body into LoginRequest struct
    │  2. c.Validate(&req) → check all validate tags (required, email, etc.)
    │  3. call service
    ▼
Service  (internal/domain/user/service.go)
    │  1. ask repository for user by email
    │  2. check password with bcrypt
    │  3. generate JWT tokens
    │  4. return response DTO
    ▼
Repository  (internal/domain/user/repository.go)
    │  run GORM query: SELECT * FROM users WHERE email = ?
    ▼
PostgreSQL database (Neon)
    │  return row
    ▼
Repository → Service → Handler → Echo → Client
    response: { "access_token": "...", "refresh_token": "..." }
```

---

## Auth Flow (JWT)

JWT = JSON Web Token. A signed string that proves who you are.

### Step 1 — Register

```
POST /api/v1/auth/register
Body: { name, email, password }

→ password is hashed with bcrypt
→ user row saved to DB
→ returns user info (no tokens yet)
```

### Step 2 — Login

```
POST /api/v1/auth/login
Body: { email, password }

→ find user by email
→ compare password with bcrypt hash
→ generate access token  (expires in 15 minutes)
→ generate refresh token (expires in 7 days)
→ return both tokens
```

### Step 3 — Access Protected Routes

```
GET /api/v1/auth/me
Header: Authorization: Bearer <access_token>

→ AuthMiddleware reads the header
→ validates the token signature using JWT_SECRET
→ if valid: sets user_id in context, calls handler
→ handler reads user_id from context, returns user info
```

### Step 4 — Refresh Access Token

Access tokens expire in 15 minutes. When they expire:

```
POST /api/v1/auth/refresh
Body: { "refresh_token": "<your_refresh_token>" }

→ validates the refresh token
→ issues a new access token
→ refresh token is still valid for 7 days
```

### Why two tokens?

| Token | Lifespan | Purpose |
|-------|----------|---------|
| Access token | 15 min | Used on every request. Short life = less risk if stolen |
| Refresh token | 7 days | Only used to get new access tokens. Never sent to API routes |

If someone steals your access token, it expires in 15 minutes. If the refresh token were used everywhere, a stolen one would give access for 7 days.
