# 🚗 SpotSync – Smart Parking & EV Charging Reservation API

Live link: https://spotsync-6576.onrender.com

SpotSync is a centralized, high-performance backend system designed to manage parking zones and handle high-demand reservations for limited EV charging spots at busy hubs like airports and shopping malls. 

Built with **Go**, **Echo**, and **GORM**, this project strictly follows **Clean Architecture** and **Domain-Driven Design (DDD)**. It features robust concurrency controls to prevent double-booking of parking spots during high-traffic scenarios.

---

## 🌟 Features & Overview

* **Role-Based Access Control (RBAC):** Distinct permissions for `admin` (management) and `driver` (booking).
* **Dynamic Capacity Calculation:** Real-time calculation of `available_spots` using database subqueries.
* **Concurrency Safe (Race Condition Handled):** Implements **GORM Transactions** alongside **PostgreSQL Row-Level Locking (`FOR UPDATE`)** to guarantee that the final EV spot is never double-booked, even if multiple drivers request it at the exact same millisecond.
* **Domain-Driven Design (DDD):** Code is highly modularized into `user`, `zone`, and `reservation` domains.
* **Centralized Error Handling:** Consistent and predictable JSON response structures for both success and error states.

---

## 🛠️ Technology Stack

* **Language:** Go (Golang) v1.22+
* **Web Framework:** Echo (`github.com/labstack/echo/v4`)
* **ORM:** GORM (`gorm.io/gorm`)
* **Database:** PostgreSQL (NeonDB / Supabase / Aiven)
* **Validation:** Go-Playground Validator v10
* **Security & Auth:** JWT (`golang-jwt/jwt/v5`) & Bcrypt hashing

---

## 🏗️ Project Architecture

The project enforces strict separation of concerns. Handlers do not interact with the database, and Repositories do not handle HTTP logic. Everything is manually wired via Dependency Injection in `main.go`.

```text
.
├── cmd/
│   └── main.go                  # Application Entrypoint & Dependency Injection (Wiring)
├── internal/
│   ├── auth/                    # JWT Token signing & validation
│   ├── config/                  # Environment variables & Database setup
│   ├── domain/                  # Isolated Business Domains
│   │   ├── user/                # Auth, Registration, Login
│   │   ├── zone/                # Zone Management & Capacity logic
│   │   └── reservation/         # Bookings & Transaction Locking (FOR UPDATE)
│   ├── httpresponse/            # Uniform API JSON wrappers
│   ├── middlewares/             # JWT Verification & Admin/Driver Role Guards
│   └── server/                  # HTTP Server & Routes setup
└── README.md