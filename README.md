# Money Transfer System

A concurrent money transfer system built in Go that allows users to transfer money between accounts with atomic updates, overdraft prevention, and a REST API.

## Features

- Transfer money between user accounts
- Atomic updates to prevent race conditions
- Overdraft prevention (users cannot send more money than they have)
- HTTP API for initiating transfers and querying user information
- Comprehensive error handling

## Locking Strategy

The system uses a mutex-based locking strategy to ensure atomic updates and prevent race conditions during money transfers:

### Coarse-Grained Locking

The implementation uses a coarse-grained locking approach in the `TransferService`:

```go
// holds a single mutex for all transfer operations
type TransferService struct {
    userRepo     repository.UserRepository
    transferRepo repository.TransferRepository
    idGenerator  *utils.IDGenerator
    mu           sync.Mutex
}
```

**Why this approach was chosen:**

1. **Simplicity**: A single mutex is easy to reason about and ensures that no race conditions can occur during transfers.
2. **Atomicity**: The entire transfer operation (checking balances, debiting sender, crediting receiver, persisting the transfer) is performed as an atomic unit.
3. **Consistency**: Since all transfers are serialized, the system will always maintain a consistent state.

**Trade-offs:**

1. **Performance**: This approach may become a bottleneck under high load since all transfers are processed serially, not in parallel.
2. **Scalability**: Limited to a single instance since the mutex only works within a single process.

### Repository-Level Concurrency Control

Each repository also implements its own mutex to protect concurrent access to the in-memory data:

```go
// has its own mutex for concurrent access
type UserRepository struct {
    users map[string]*model.User
    mu    sync.RWMutex
}
```

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/IskenT/money-transfer.git
   cd money-transfer
   ```

2. Build the application:
   ```bash
   make build
   ```

3. Run the application:
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080` with Swagger documentation at `http://localhost:8080/swagger/index.html`.

### Running Tests

```bash
make test
```

## API Endpoints

- `POST /api/transfers` - Create a new transfer
- `GET /api/transfers` - List all transfers
- `GET /api/transfers/{id}` - Get transfer details
- `GET /api/users` - List all users with their balances
- `GET /api/users/{id}` - Get user details by ID

## Initial Account Balances

- Mark: $100.00
- Jane: $50.00
- Adam: $0.00