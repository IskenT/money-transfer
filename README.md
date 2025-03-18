# Money Transfer System

A robust money transfer system built in Go that allows users to transfer money between accounts with atomic updates, overdraft prevention, and a REST API. This version uses PostgreSQL for persistent storage and implements the Transactional Outbox Pattern.

## Key Features

- Transfer money between user accounts
- Atomic database transactions with proper isolation levels
- Row-level locking with SELECT FOR UPDATE to prevent race conditions
- Outbox pattern for reliable event publishing
- PostgreSQL database for persistent storage
- RESTful API with Swagger documentation
- Automatic database migration management

## Architecture Overview

### Database Design

The system uses PostgreSQL with the following schema:
- `users` table for storing user information and balances
- `transactions` table for individual debit and credit transactions
- `transfers` table for tracking money transfers between users
- `outbox_events` table for the transactional outbox pattern

### Concurrency Control

The system implements robust concurrency control using:

1. **Database Transactions**: All transfer operations occur within a database transaction to ensure atomicity.
2. **Row-Level Locking**: Using `SELECT FOR UPDATE` to lock rows during balance updates.
3. **Isolation Level**: Transactions use REPEATABLE READ isolation to prevent dirty, non-repeatable, and phantom reads.

### Transactional Outbox Pattern

The system uses the Transactional Outbox Pattern to reliably publish events after a successful transfer:

1. Transfer data and related events are written to the database in a single transaction
2. A background processor periodically polls the outbox table for unprocessed events
3. Events are processed and marked as completed

## Getting Started

### Prerequisites

- Go 1.19+
- Podman (for PostgreSQL container)

### Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/IskenT/money-transfer.git
   cd money-transfer
   ```

2. Start the application (sets up database and runs migrations automatically):
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080` with Swagger documentation at `http://localhost:8080/swagger/index.html`.

### Additional Commands

- **Build the application**:
  ```bash
  make build
  ```

- **Generate Swagger documentation**:
  ```bash
  make docs
  ```

- **Run tests**:
  ```bash
  make test
  ```

- **Stop database**:
  ```bash
  make db-stop
  ```

- **Reset everything** (useful for a clean slate):
  ```bash
  make reset
  ```

## API Endpoints

- `POST /api/transfers` - Create a new transfer
- `GET /api/transfers` - List all transfers
- `GET /api/transfers/{id}` - Get transfer details by ID
- `GET /api/users` - List all users with their balances
- `GET /api/users/{id}` - Get user details by ID

## Initial Account Balances

- Mark: $100.00
- Jane: $50.00
- Adam: $0.00

## Example API Requests

### Create a transfer

```bash
curl -X POST http://localhost:8080/api/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "from_user_id": "1",
    "to_user_id": "2",
    "amount": 2000
  }'
```

### List all users

```bash
curl -X GET http://localhost:8080/api/users
```

## Project Structure

```
money-transfer/
├── cmd/
│   ├── migrate/       # Database migration tool
│   └── server/        # Main application entry point
├── docs/              # Swagger documentation
├── internal/
│   ├── app/
│   │   ├── processor/ # Background processors (outbox)
│   │   └── service/   # Business logic services
│   ├── application/   # Application setup
│   ├── config/        # Configuration management
│   ├── domain/
│   │   ├── model/     # Domain models
│   │   └── repository/# Repository interfaces
│   └── infra/
│       ├── database/  # Database connection and transaction management
│       ├── http/      # HTTP handlers, routers, and models
│       └── repository/# Repository implementations
├── migrations/        # SQL migration files
├── docker-compose.yml  # Podman container configuration
└── Makefile           # Build and run commands
```



![alt text](image.png)