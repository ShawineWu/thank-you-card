.PHONY: help backend-dev backend-build backend-run frontend-dev frontend-build frontend-install db-init db-reset clean

# Variables
BACKEND_DIR=backend
FRONTEND_DIR=frontend
DB_NAME=thankyoucard
DB_USER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=123456

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Backend targets
backend-dev: ## Start backend development server
	@echo "Starting backend server..."
	@cd $(BACKEND_DIR) && \
		DATABASE_DSN="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		PORT=8056 \
		MOCK_USER_ID=1 \
		go run ./cmd/server

backend-build: ## Build backend binary
	@echo "Building backend..."
	@cd $(BACKEND_DIR) && go build -o bin/server ./cmd/server

backend-run: backend-build ## Build and run backend
	@echo "Running backend server..."
	@cd $(BACKEND_DIR) && \
		DATABASE_DSN="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		PORT=8056 \
		MOCK_USER_ID=1 \
		./bin/server

# Frontend targets
frontend-install: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && npm install

frontend-dev: ## Start frontend development server
	@echo "Starting frontend development server..."
	@cd $(FRONTEND_DIR) && npm run dev

frontend-build: ## Build frontend for production
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && npm run build

# Database targets
db-init: ## Initialize database (create if not exists)
	@echo "Initializing database..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 || \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "Database '$(DB_NAME)' is ready"

db-reset: ## Reset database (WARNING: drops and recreates database)
	@echo "WARNING: This will drop the database!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"; \
		echo "Database reset complete"; \
	fi

# Development workflow
dev: db-init backend-dev ## Initialize DB and start backend (in one terminal)

dev-full: db-init ## Start both backend and frontend (requires two terminals)
	@echo "Starting full development environment..."
	@echo "Backend will run in this terminal. Open another terminal and run 'make frontend-dev'"
	@$(MAKE) backend-dev

# Setup for new developers
setup: db-init frontend-install ## Complete setup for new developers
	@echo ""
	@echo "Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Start backend: make backend-dev"
	@echo "2. Start frontend (in another terminal): make frontend-dev"
	@echo "3. Open http://localhost:3000 in your browser"

# Clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BACKEND_DIR)/bin
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf $(FRONTEND_DIR)/node_modules/.vite
	@echo "Clean complete"

# Quick start (assumes DB is already running)
start: db-init ## Quick start: init DB and show instructions
	@echo ""
	@echo "Database initialized. Now you can:"
	@echo "  Terminal 1: make backend-dev"
	@echo "  Terminal 2: make frontend-dev"
