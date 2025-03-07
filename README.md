# Money Transfer Service

A simple and reliable money transfer service built with Go. This service provides basic money transfers between accounts while maintaining data consistency and preventing race conditions.

## 🚀 Key Features

- ⚡️ Atomic money transfers
- 🔒 Race condition prevention using SERIALIZABLE isolation
- 💰 Overdraft protection
- ✅ Basic test coverage

## 🛠 Tech Stack

- **Go** - core language
- **PostgreSQL** - database
- **Docker** - containerization
- **Viper** - configuration
- **Testify** - testing
- **Swagger** - API documentation
- **golangci-lint** - code quality tool

## 🏃‍♂️ Quick Start

### Using Docker

```bash
# Clone the repository
git clone https://github.com/ZakDziu/money-transfer
cd money-transfer

# Start the service
docker-compose up -d

# Check if it's working
curl http://localhost:8080/api/v1/balance/Mark
```

### Local Development

```bash
# Install dependencies
go mod download

# Start PostgreSQL
docker-compose up -d postgres

# Run the service
go run cmd/server/main.go
```

## 📡 API

### Transfer Money

```bash
POST /api/v1/transfer
Content-Type: application/json

{
    "from": "Mark",
    "to": "Jane",
    "amount": 50.0
}
```

### Check Balance

```bash
GET /api/v1/balance/{account}
```

### API Documentation
Full API documentation is available via Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

## 🧪 Testing

```bash
# Run all tests with test database setup
chmod +x ./scripts/run-tests.sh
./scripts/run-tests.sh

# Run tests with coverage
GO_ENV=test go test -v -cover ./...

# Run tests with race detector
GO_ENV=test go test -v -race ./...

# Run linter
golangci-lint run
```

## 📁 Project Structure

```
.
├── cmd/                  # Application entrypoints
│   └── server/          # HTTP server
├── config/              # Configuration
├── .golangci.yml       # Linter configuration
├── internal/            # Internal code
│   ├── api/            # API layer
│   │   ├── docs/       # Swagger documentation
│   │   ├── handlers/   # Request handlers
│   │   └── router/     # Routing setup
│   ├── domain/         # Business models and errors
│   ├── service/        # Business logic
│   └── storage/        # Data storage
└── docker-compose.yml  # Docker configuration
```

## ⚙️ Configuration

The service uses environment files for configuration:

### Main Configuration (`.env`)
```env
# Server Configuration
SERVER_PORT=8080              # HTTP server port
SERVER_READ_TIMEOUT=5s        # Maximum duration for reading request
SERVER_WRITE_TIMEOUT=10s      # Maximum duration for writing response
SERVER_IDLE_TIMEOUT=15s       # Maximum duration for idle connections

# Database Configuration
DB_HOST=postgres             # PostgreSQL host
DB_PORT=5432                # PostgreSQL port
DB_USER=postgres            # Database user
DB_PASSWORD=postgres        # Database password
DB_NAME=money_transfer      # Database name
DB_SSLMODE=disable         # SSL mode for database connection
```

### Test Configuration (`.env.test`)
```env
# Database Configuration
DB_HOST=localhost          # Local PostgreSQL for tests
DB_PORT=5433              # Different port to avoid conflicts
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=money_transfer_test
DB_SSLMODE=disable
```

## 🎯 Design Decisions

### Transaction Management
- Uses PostgreSQL's SERIALIZABLE isolation level
- Single-phase commit for atomic operations
- Row-level locking to prevent deadlocks

### Error Handling
- Domain-specific error types:
  - Account not found
  - Insufficient funds
  - Invalid amount
  - Same account transfer

### Code Quality
- Strict linting rules with golangci-lint
- Consistent code style
- Static code analysis

## 👥 Authors

- **Zakhar Dziuniak** - *Initial work* - [ZakDziu](https://github.com/ZakDziu)

## 🗺 Roadmap

- [ ] Multi-currency support
- [ ] Transaction scheduling
- [ ] WebSocket notifications
- [ ] Account statements
- [ ] Batch transfers
- [ ] API versioning

---

[⬆ back to top](#money-transfer-service)