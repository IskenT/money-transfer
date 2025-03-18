# Money Transfer System - Package Structure

This document outlines the recommended package structure for the Money Transfer System using Hexagonal Architecture and Domain-Driven Design principles.

## Project Layout

```
money-transfer/
├── api/                     # API Documentation
│   ├── swagger.json
│   └── swagger.yaml
├── cmd/                     # Entry points for the application
│   └── server/
│       └── main.go          # Main entry point for the server
├── docs/                    # Swagger generated documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/                # Private application code
│   ├── domain/              # Domain Layer - Core business logic
│   │   ├── model/           # Domain models
│   │   │   ├── user.go
│   │   │   ├── transaction.go
│   │   │   ├── transfer.go
│   │   │   └── errors.go
│   │   └── repository/      # Repository interfaces
│   │       ├── user_repository.go
│   │       └── transfer_repository.go
│   ├── app/                 # Application Layer - Use cases
│   │   └── service/
│   │       └── transfer_service.go
│   └── infra/               # Infrastructure Layer
│       ├── repository/      # Repository implementations
│       │   └── memory/      # In-memory storage implementation
│       │       ├── user_repository.go
│       │       └── transfer_repository.go
│       └── http/            # HTTP-related code
│           ├── handler/     # HTTP handlers
│           │   ├── transfer_handler.go
│           │   └── user_handler.go
│           ├── middleware/  # HTTP middleware
│           │   └── cors.go
│           └── rest/        # REST API components
│               └── server.go
├── pkg/                     # Public library code
│   └── utils/               # Shared utilities
├── main.go                  # Main entry point
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
└── README.md                # Project documentation
```

## Package Descriptions

### Domain Layer (`internal/domain/`)

The core business logic and entities, independent of external concerns.

- **`model/`**: Contains the domain entities and value objects
  - `user.go`: User entity with ID, name, and balance
  - `transaction.go`: Transaction entity for debit/credit operations
  - `transfer.go`: Transfer entity linking debit and credit transactions
  - `errors.go`: Domain-specific errors

- **`repository/`**: Interfaces that define how to access domain entities
  - `user_repository.go`: Interface for user persistence operations
  - `transfer_repository.go`: Interface for transfer persistence operations

### Application Layer (`internal/app/`)

Use cases and services that orchestrate the domain objects.

- **`service/`**: Application services that implement use cases
  - `transfer_service.go`: Service for creating and managing transfers

### Infrastructure Layer (`internal/infra/`)

Implementation of interfaces defined in the domain layer.

- **`repository/memory/`**: In-memory implementations of repository interfaces
  - `user_repository.go`: In-memory user storage
  - `transfer_repository.go`: In-memory transfer storage

- **`http/`**: HTTP-related code
  - `handler/`: HTTP handlers for API endpoints
  - `middleware/`: HTTP middleware components
  - `rest/`: REST API server setup

### Entry Points (`cmd/`)

Entry points to the application.

- **`server/`**: HTTP server entry point
  - `main.go`: Configures and starts the HTTP server

### API Documentation (`api/` and `docs/`)

Swagger documentation for the REST API.