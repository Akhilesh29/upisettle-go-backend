## UpiSettle (backend)

- Go backend to reconcile UPI and cash payments against store orders, so merchants get a clear end-of-day settlement view.
- Exposes a simple JSON API used by mobile/web clients.

### High-level flow

```mermaid
flowchart LR
    C[Customer] -->|UPI / Cash payment| M[Merchant Store]
    M -->|Creates order| OS[Order Service]
    UPI[Bank / UPI App] -->|SMS / webhook data| MP[Mobile Parser App]
    MP -->|Parsed UPI payment event| PS[Payment Service]
    OS -->|Orders & cash payments| DB[(PostgreSQL)]
    PS -->|UPI payments| DB
    RECO[Reconciliation Service] -->|Reads orders & payments| DB
    RECO --> SUM[Daily Summary & Exceptions]
    SUM --> OWN[Merchant Owner (mobile/web)]
```

---

# UpiSettle Backend (Go)

UpiSettle is a backend service that helps small offline merchants (kirana shops, retail stores) **reconcile UPI and cash payments against their orders** and get a clear end‑of‑day settlement view.

This repository contains the **backend-only** implementation written in Go, designed as a **modular monolith** that is easy to extend and scale.

---

## Features (MVP)

- **Merchant & stores**
  - Register a merchant owner.
  - Create and list stores for a merchant.
- **Auth**
  - JWT-based authentication for API access.
- **Orders**
  - Create orders per store.
  - List orders for a given day.
- **Payments**
  - Ingest parsed UPI payment events (from mobile app SMS parser).
  - Record manual cash payments against orders.
- **Reconciliation**
  - Match orders and UPI payments by amount for a given day.
  - Create exceptions for unmatched orders/payments or ambiguous matches.
- **Reporting**
  - Per-store daily summary (sales, UPI vs cash totals, matched vs unmatched, exceptions).
  - List exceptions for a given day.

---

## Architecture Overview

High-level layers:

- `cmd/api`: application entrypoint (`main.go`).
- `internal/config`: environment-based configuration (port, DB URL, JWT secret).
- `internal/logger`: simple structured logging wrapper.
- `internal/storage`: database connection (PostgreSQL via GORM).
- `internal/http`: HTTP server setup with Gin, global middlewares, and route wiring.
- Domain modules:
  - `internal/auth`: users, registration, login, JWT middleware.
  - `internal/merchant`: merchants and stores.
  - `internal/order`: orders and basic listing.
  - `internal/payment`: payment ingestion (UPI & cash).
  - `internal/matching`: reconciliation engine and models (`matches`, `exceptions`).
  - `internal/reporting`: daily summaries and exception listings.
- `migrations`: SQL migrations for the relational schema.

The API service is **stateless**; all data is stored in PostgreSQL. This allows horizontal scaling by running multiple instances behind a load balancer.

---

## Tech Stack

- **Language**: Go (>= 1.21)
- **Web framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/) with the Postgres driver
- **Auth**: JWT (HMAC) + bcrypt password hashing

---

## Quickstart

### Prerequisites

- Go **1.21+** installed and available on your `PATH`.
- Docker & Docker Compose (for running PostgreSQL locally).

### 1. Clone and navigate

```bash
cd "c:/Users/HP/Desktop/New folder"
cd UpiSettle
```

### 2. Start PostgreSQL

```bash
docker-compose up -d db
```

This starts a Postgres instance on `localhost:5432` with:

- user: `upisettle`
- password: `upisettle`
- database: `upisettle`

### 3. Configure environment

Copy the example env file and adjust values if needed:

```bash
cp .env.example .env
```

Key variables:

- `APP_ENV` – `development` / `production`.
- `PORT` – HTTP server port (default `8080`).
- `DATABASE_URL` – Postgres DSN.
- `JWT_SECRET` – secret key for signing JWTs (change this!).

### 4. Run database migrations

Migrations are defined in the `migrations` folder (e.g. `0001_init_schema.up.sql`).

You can apply them using [`golang-migrate`](https://github.com/golang-migrate/migrate) or your preferred migration tool. Example command (adjust paths as needed):

```bash
migrate -path ./migrations -database "postgres://upisettle:upisettle@localhost:5432/upisettle?sslmode=disable" up
```

### 5. Install dependencies and run

From the `UpiSettle` directory:

```bash
go mod tidy
go run ./cmd/api
```

The server will start on `http://localhost:8080` (or your configured `PORT`).

Health check:

```bash
curl http://localhost:8080/healthz
```

---

## API Overview (MVP)

Base path: `/api/v1`

All responses are JSON. For authenticated endpoints, include:

```http
Authorization: Bearer <token>
```

### Auth

- **POST** `/api/v1/auth/register`

  Register a merchant owner and create the merchant.

  Request body:

  ```json
  {
    "name": "Akhilesh",
    "email": "owner@example.com",
    "phone": "9999999999",
    "password": "secret123",
    "merchant_name": "Akhilesh Kirana"
  }
  ```

  Response:

  ```json
  { "token": "<JWT_TOKEN>" }
  ```

- **POST** `/api/v1/auth/login`

  Request:

  ```json
  {
    "email": "owner@example.com",
    "password": "secret123"
  }
  ```

  Response:

  ```json
  { "token": "<JWT_TOKEN>" }
  ```

### Merchant & Stores (Authenticated)

- **POST** `/api/v1/stores`

  Create a store for the authenticated merchant.

  ```json
  {
    "name": "Main Store",
    "address": "Some street, City"
  }
  ```

- **GET** `/api/v1/stores`

  List stores for the authenticated merchant.

### Orders (Authenticated)

- **POST** `/api/v1/stores/{storeId}/orders`

  Create an order for a store.

  ```json
  {
    "amount": 15000,
    "external_ref": "ORDER-101"
  }
  ```

  Amount is in **paise** (smallest unit), so `15000` = ₹150.00.

- **GET** `/api/v1/stores/{storeId}/orders?date=YYYY-MM-DD`

  List orders for a specific date.

### Payments (Authenticated)

- **POST** `/api/v1/stores/{storeId}/payments`

  Ingest a payment event (typically UPI) for a store.

  ```json
  {
    "channel": "UPI",
    "amount": 15000,
    "time": "2026-02-13T17:10:00Z",
    "upi_ref": "123456",
    "payer_vpa": "user@upi",
    "payer_name": "Some User"
  }
  ```

- **POST** `/api/v1/stores/{storeId}/cash-payments`

  Record a cash payment against an order and mark it as paid by cash.

  ```json
  {
    "order_id": 101,
    "amount": 15000
  }
  ```

### Reconciliation (Authenticated)

- **POST** `/api/v1/stores/{storeId}/reconcile?date=YYYY-MM-DD`

  Run reconciliation for a given store and date. Matches orders with payments by amount, creates match records and exceptions.

  Example response:

  ```json
  {
    "matched_orders": 10,
    "unmatched_orders": 2,
    "unmatched_payments": 1
  }
  ```

### Reporting (Authenticated)

- **GET** `/api/v1/stores/{storeId}/summary?date=YYYY-MM-DD`

  Returns a daily summary:

  ```json
  {
    "date": "2026-02-13",
    "total_orders": 12,
    "total_sales_amount": 180000,
    "upi_total_amount": 120000,
    "cash_total_amount": 60000,
    "matched_orders": 10,
    "unmatched_orders": 2,
    "exceptions_count": 3,
    "exceptions_amount": 30000
  }
  ```

- **GET** `/api/v1/stores/{storeId}/exceptions?date=YYYY-MM-DD`

  Returns a list of reconciliation exceptions for the day.

---

## Development Guidelines

- Keep the **layering** clean:
  - HTTP handlers (`internal/*/http.go`) should be thin: parse/validate input, call services, format responses.
  - Services (`internal/*/service.go`) should contain business logic and talk to the DB via GORM.
  - Models live close to their domain (`internal/auth`, `internal/merchant`, etc.).
- When adding a new module:
  1. Define models and services under `internal/<module>/`.
  2. Expose HTTP handlers via a `RegisterHTTP` function.
  3. Wire the module in `internal/http/server.go` under the appropriate route group.
- Use **paise** (int64) for monetary values to avoid floating-point issues.

---

## Next Steps / Possible Extensions

- More advanced matching heuristics (time windows, partial payments, customer hints).
- Background job for automatic reconciliation instead of only manual trigger.
- Better role management (owner vs staff).
- Audit logging and idempotency for payment ingestion.

