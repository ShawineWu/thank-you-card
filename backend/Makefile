.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down deps install

# Variables
BINARY_NAME=thank-you-card-backend
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
GO=go
DOCKER_COMPOSE=docker-compose

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run
run: ## Run the application (requires .env file)
	@echo "Starting server..."
	@$(GO) run $(MAIN_PATH)

# Install dependencies
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "Dependencies updated"

# Test
test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Database operations
db-init: ## Initialize database using init.sql (requires PostgreSQL client)
	@echo "Initializing database..."
	@PGPASSWORD=123456 psql -h localhost -p 5432 -U postgres -f script/init.sql
	@echo "Database initialized"

db-reset: ## Reset database (WARNING: This will drop and recreate the database)
	@echo "WARNING: This will drop the database!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		PGPASSWORD=123456 psql -h localhost -p 5432 -U postgres -d postgres -c "DROP DATABASE IF EXISTS thankyoucard;"; \
		PGPASSWORD=123456 psql -h localhost -p 5432 -U postgres -f script/init.sql; \
		echo "Database reset complete"; \
	fi

# Docker operations
docker-up: ## Start PostgreSQL using Docker Compose
	@echo "Starting PostgreSQL container..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database is ready"

docker-down: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	@$(DOCKER_COMPOSE) down
	@echo "Container stopped"

docker-logs: ## View Docker Compose logs
	@$(DOCKER_COMPOSE) logs -f

# Development workflow
dev: docker-up run ## Start database and run the application

dev-clean: docker-down clean ## Stop containers and clean build artifacts

# Linting and formatting
fmt: ## Format Go code
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...

lint: fmt vet ## Run all linters

# Quick setup for new developers
setup: deps docker-up ## Setup development environment
	@echo "Setup complete!"
	@echo "Next steps:"
	@echo "1. Copy .env.example to .env and configure it"
	@echo "2. Run 'make run' to start the server"
