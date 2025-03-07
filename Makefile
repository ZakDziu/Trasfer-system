# Variables
APP_NAME = money-transfer
MAIN_PATH = cmd/server/main.go
BINARY_PATH = bin/server

# Go variables
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./internal/api/docs/*")

# Docker variables
DOCKER_COMPOSE = docker-compose
DOCKER_IMAGE = money-transfer

.PHONY: all build run test test-coverage lint clean help docker-build docker-up docker-down install-deps generate-swagger

# Main commands
all: install-deps lint test build ## Run all main tasks

build: ## Build the application
	go build -o $(BINARY_PATH) $(MAIN_PATH)

run: ## Run the application
	go run $(MAIN_PATH)

# Test commands
test-script: ## Run tests using script with database setup
	chmod +x ./scripts/run-tests.sh
	./scripts/run-tests.sh

test: ## Run tests
	GO_ENV=test go test -v ./...

test-coverage: ## Run tests with coverage
	GO_ENV=test go test -v -cover ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-race: ## Run tests with race detector
	GO_ENV=test go test -v -race ./...

# Lint commands
lint: ## Run linter
	chmod +x ./scripts/lint.sh
	./scripts/lint.sh

lint-fix: ## Run linter with auto-fix
	golangci-lint run --fix

# Docker commands
docker-build: ## Build docker image
	$(DOCKER_COMPOSE) build

docker-up: ## Start docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop docker containers
	$(DOCKER_COMPOSE) down

docker-logs: ## Show docker logs
	$(DOCKER_COMPOSE) logs -f

# Database commands
db-up: ## Start database
	$(DOCKER_COMPOSE) up -d postgres

db-down: ## Stop database
	$(DOCKER_COMPOSE) down postgres

db-test-up: ## Start test database
	$(DOCKER_COMPOSE) up -d postgres_test

# Development tools
install-deps: ## Install development dependencies
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.55.2
	chmod +x ./scripts/lint.sh

generate-swagger: ## Generate Swagger documentation
	swag init -g $(MAIN_PATH) -o internal/api/docs

# Utility commands
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf .golangci-lint-cache/

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-tidy: ## Tidy go modules
	go mod tidy

# Development workflow
dev: ## Run application in development mode
	air -c .air.toml

# Help command
help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help 