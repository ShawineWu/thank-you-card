.PHONY: up down restart logs install dev-backend dev-frontend help

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install dependencies for all projects"
	@echo "  make up            - Start database (and other docker services) in background"
	@echo "  make down          - Stop all docker services"
	@echo "  make restart       - Restart all docker services"
	@echo "  make logs          - View database logs"
	@echo "  make dev-backend   - Run backend services (requires 'make up' first)"
	@echo "  make dev-frontend  - Run frontend services"
	@echo "  make clean         - Remove project temporary files and build artifacts"

install:
	@echo "Installing Backend Dependencies..."
	cd src/backend && go mod download
	@echo "Installing Frontend (Teams App) Dependencies..."
	cd src/frontend/teams-app && pnpm install
	@echo "Installing Frontend (Web Dashboard) Dependencies..."
	cd src/frontend/web-dashboard && pnpm install

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

dev-backend:
	@echo "Starting Backend Services..."
	@echo "Ensure you have run 'make up' to start the database."
	cd src/backend && go run cmd/card-service/main.go & \
	go run cmd/analytics-service/main.go & \
	wait

dev-frontend:
	@echo "Starting Frontend Services..."
	(cd src/frontend/teams-app && pnpm dev) & \
	(cd src/frontend/web-dashboard && pnpm dev) & \
	wait

clean:
	rm -rf src/frontend/teams-app/node_modules
	rm -rf src/frontend/teams-app/dist
	rm -rf src/frontend/web-dashboard/node_modules
	rm -rf src/frontend/web-dashboard/dist
