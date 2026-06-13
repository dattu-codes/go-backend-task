# Go User Management REST API

A premium, production-ready REST API built in Go for managing users, structured using Clean Architecture principles (Handler, Service, Repository, Database).


---

## ✨ Features

- **CRUD Operations**: Support for Create, Read, Update, and Delete endpoints for users.
- **Dynamic Age Calculation**: Dynamically computes exact user age in years using calendar boundaries, month/day offset checks, and Feb 29 leap year cases.
- **SQLC Database Client**: Generates 100% type-safe SQL query bindings on top of standard schema SQL definitions.
- **pgx/v5 Driver & Connection Pooling**: Utilizes connection pools via `pgxpool` for concurrent safety.
- **Uber Zap Structured Logging**: Clean, structured JSON logs containing timestamp, log level, caller, method, latency, and request-id.
- **Custom Request-ID Middleware**: Propagates or generates a cryptographically secure UUID (`crypto/rand`) as `X-Request-ID` across response headers and log scopes.
- **Response Latency Logging**: Measures and outputs route execution duration.
- **List Pagination**: High-performance pagination with `page` and `limit` query parameters.
- **Docker Compose Orchestration**: Containerizes the application and PostgreSQL database with standard health checking.

---

## 🛠️ Technology Stack & Decisions

- **Language**: Go 1.22+
- **HTTP Routing**: GoFiber v2 (Zero-allocation routing built on `fasthttp`)
- **DB Driver**: jackc/pgx/v5 (modern PostgreSQL client)
- **DB Compiler**: SQLC (SQL-first compiler producing compiled Go helper structures)
- **Validation**: go-playground/validator/v10
- **Logging**: Uber Zap (JSON structured logger)

---

## 📁 Repository Structure

```text
├── cmd/
│   └── server/
│       └── main.go           # Application entrypoint & connection setup
├── config/
│   └── config.go         # Config loader reading env variables with defaults
├── db/
│   ├── migrations/       # Schema change files (up/down migrations)
│   ├── query.sql         # SQL query targets compiled by SQLC
│   ├── schema.sql        # Table structures compiled by SQLC
│   └── sqlc/             # Generated database methods (db.go, models.go, query.sql.go)
├── internal/
│   ├── handler/          # User HTTP controllers
│   ├── repository/       # DB interfaces separating controllers from SQL engines
│   ├── service/          # Core domain logic (e.g., dynamic age calculation)
│   ├── routes/           # Routing configuration & route grouping
│   ├── middleware/       # Custom middlewares (RequestId injection, Logger)
│   ├── models/           # Struct definitions (Request / Response validation models)
│   └── logger/           # Configured Uber Zap logger
├── Dockerfile            # Lightweight, multi-stage compilation image
├── docker-compose.yml    # Postgres database container + web app container
└── sqlc.yaml             # SQLC compiler configuration mapping
```

---

## 🚀 How to Run the App

### Option A: Using Docker Compose (Recommended)

Docker Compose configures a PostgreSQL container and the Go API container, establishing network connections automatically with standard health checks.

1. Build and boot the stack:
   ```bash
   docker-compose up --build
   ```
2. The database schema will be initialized inside PostgreSQL automatically, and the API will boot up on port `3000`.

### Option B: Running Natively

1. **Start PostgreSQL**: Make sure you have PostgreSQL running and create a database named `userdb`.
2. **Apply Schema**: Run the schema SQL inside `db/schema.sql` on your database:
   ```bash
   psql -U postgres -d userdb -f db/schema.sql
   ```
3. **Configure Environment**: Set the database URL:
   ```bash
   # Windows (PowerShell)
   $env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"
   $env:PORT="3000"
   
   # Linux/macOS
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"
   export PORT="3000"
   ```
4. **Run Server**:
   ```bash
   go run cmd/server/main.go
   ```

---

## 🧪 Running Tests

To execute the unit tests for age calculation and calendar logic:
```bash
go test -v ./internal/service/...
```

---

## 📝 API Endpoints

All responses include the calculated `age` dynamically, along with a custom `X-Request-ID` response header.

### 1. Create User
- **Method**: `POST`
- **Path**: `/users`
- **Headers**: `Content-Type: application/json`
- **Payload**:
  ```json
  {
    "name": "Alice Smith",
    "dob": "1995-10-15"
  }
  ```
- **Response** (HTTP 201 Created):
  ```json
  {
    "id": 1,
    "name": "Alice Smith",
    "dob": "1995-10-15",
    "age": 30
  }
  ```

### 2. Get User
- **Method**: `GET`
- **Path**: `/users/:id`
- **Response** (HTTP 200 OK):
  ```json
  {
    "id": 1,
    "name": "Alice Smith",
    "dob": "1995-10-15",
    "age": 30
  }
  ```
- **Error Response** (HTTP 404 Not Found):
  ```json
  {
    "error": "User not found"
  }
  ```

### 3. Update User
- **Method**: `PUT`
- **Path**: `/users/:id`
- **Payload**:
  ```json
  {
    "name": "Alice Cooper",
    "dob": "1995-10-20"
  }
  ```
- **Response** (HTTP 200 OK):
  ```json
  {
    "id": 1,
    "name": "Alice Cooper",
    "dob": "1995-10-20",
    "age": 30
  }
  ```

### 4. Delete User
- **Method**: `DELETE`
- **Path**: `/users/:id`
- **Response** (HTTP 204 No Content)

### 5. List Users (with Pagination)
- **Method**: `GET`
- **Path**: `/users?page=1&limit=5`
- **Query Params**:
  - `page` (optional, default: 1): Target page.
  - `limit` (optional, default: 10): Records per page.
- **Response** (HTTP 200 OK):
  ```json
  [
    {
      "id": 1,
      "name": "Alice Cooper",
      "dob": "1995-10-20",
      "age": 30
    }
  ]
  ```
